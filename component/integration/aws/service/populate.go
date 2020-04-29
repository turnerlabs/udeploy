package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Populate ...
func Populate(instances map[string]app.Instance) (map[string]app.Instance, error) {

	populated := map[string]app.Instance{}

	for name, instance := range instances {

		session := session.New()

		config.Merge([]string{instance.Role}, session)

		svc := ecs.New(session)

		resourceIds := []string{}
		for _, i := range instances {
			resourceIds = append(resourceIds, fmt.Sprintf("service/%s/%s", i.Cluster, i.Service))
		}

		ascv := applicationautoscaling.New(session)

		ao, err := ascv.DescribeScalableTargets(&applicationautoscaling.DescribeScalableTargetsInput{
			ServiceNamespace: aws.String("ecs"),
			ResourceIds:      aws.StringSlice(resourceIds),
		})
		if err != nil {
			return instances, err
		}

		i, state, err := populateInst(instance, ao.ScalableTargets, svc)
		if err != nil {
			state.SetError(err)
		}

		i.SetState(state)

		populated[name] = i
	}

	return populated, nil
}

func populateInst(i app.Instance, scalableTargets []*applicationautoscaling.ScalableTarget, svc *ecs.ECS) (app.Instance, app.State, error) {
	i.Task.Definition = app.NewDefinition(i.Task.Family)

	state := app.NewState()

	td, svcs, err := getServiceInfo(i, svc)
	if err != nil {
		return i, state, err
	}

	i.Task.Definition, err = app.DefinitionFrom(td, i.Task.ImageTagEx)
	if err != nil {
		state.SetError(err)
	}

	state.Version = i.Task.Definition.Version.Full()

	stoppedTasks, err := getTaskDetails(svc, i, []*ecs.Task{}, "STOPPED", "")
	if err != nil {
		return i, state, err
	}

	runningTasks, err := getTaskDetails(svc, i, []*ecs.Task{}, "RUNNING", "")
	if err != nil {
		return i, state, err
	}

	tasks := append(runningTasks, stoppedTasks...)

	if err := checkError(svcs, stoppedTasks, app.FailedTaskExpiration*time.Minute); err != nil {
		state.SetError(err)
		state.SetPending()
	} else if isPending(svcs) {
		state.SetPending()
	} else if isStopped(svcs) {
		state.SetStopped()
	} else {
		state.SetRunning()
	}

	i.Task.DesiredCount = *svcs.DesiredCount

	for _, t := range scalableTargets {
		if *t.ResourceId == fmt.Sprintf("service/%s/%s", i.Cluster, i.Service) {
			i.Task.DesiredCount = *t.MinCapacity
			i.AutoScale = true
		}
	}

	i.Task.TasksInfo, err = task.GetTasksInfo(i, svc, tasks)
	if err != nil {
		return i, state, err
	}

	region, err := getRegion(*td.TaskDefinitionArn)
	if err != nil {
		return i, state, err
	}

	linkName := "logs"
	if missingLink(linkName, i.Links) {
		i.Links = append(i.Links, app.Link{
			Generated:   true,
			Description: "AWS Console Service Logs",
			Name:        linkName,
			URL: fmt.Sprintf("https://console.aws.amazon.com/ecs/home?region=%s#/clusters/%s/services/%s/logs",
				region, i.Cluster, i.Service),
		})
	}

	return i, state, nil
}

func missingLink(name string, links []app.Link) bool {
	for _, l := range links {
		if l.Generated && l.Name == name {
			return false
		}
	}

	return true
}

func getRegion(arn string) (string, error) {
	tag := regexp.MustCompile("([a-z]{2}-[a-z]*-[0-9]{1})")

	matches := tag.FindStringSubmatch(arn)
	if matches == nil {
		return "", errors.New("failed to get region")
	}

	if len(matches) >= 2 && len(matches[1]) > 0 {
		return matches[1], nil
	}

	return "", errors.New("failed to get region")
}

func checkError(svcs *ecs.Service, tasks []*ecs.Task, errorExpiration time.Duration) error {

	if *svcs.DesiredCount == 0 {
		return nil
	}

	if _, err := getServiceError(tasks, errorExpiration); err != nil {
		return app.InstanceError{Problem: err.Error()}
	}

	return nil
}

func getTaskDetails(svc *ecs.ECS, inst app.Instance, tasks []*ecs.Task, status, nextToken string) ([]*ecs.Task, error) {
	input := &ecs.ListTasksInput{
		Cluster:       aws.String(inst.Cluster),
		ServiceName:   aws.String(inst.Service),
		DesiredStatus: aws.String(status),
	}

	if len(nextToken) > 0 {
		input.SetNextToken(nextToken)
	}

	stoppedTasks, err := svc.ListTasks(input)
	if err != nil {
		return nil, err
	}

	if len(stoppedTasks.TaskArns) > 0 {
		taskDetails, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: aws.String(inst.Cluster),
			Tasks:   stoppedTasks.TaskArns,
		})
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, taskDetails.Tasks...)
	}

	if stoppedTasks.NextToken == nil || len(*stoppedTasks.NextToken) == 0 {
		return tasks, nil
	}

	return getTaskDetails(svc, inst, tasks, status, nextToken)
}

func getServiceError(tasks []*ecs.Task, expiration time.Duration) (int, error) {
	var reason error
	count := 0

	for _, t := range tasks {
		if t.StopCode != nil && t.StoppedReason != nil {
			if *t.StopCode != ecs.TaskStopCodeUserInitiated {

				if time.Now().Sub(*t.ExecutionStoppedAt) < expiration {
					reason = errors.New(*t.StoppedReason)
					count++
				}
			}
		}
	}

	return count, reason
}

func isPending(svc *ecs.Service) bool {
	return (len(svc.Deployments) > 1 && *svc.DesiredCount > 0) || *svc.DesiredCount > 0 && *svc.RunningCount == 0
}

func isStopped(svc *ecs.Service) bool {
	return *svc.RunningCount == 0 && *svc.PendingCount == 0
}

func getServiceInfo(instance app.Instance, svc *ecs.ECS) (*ecs.TaskDefinition, *ecs.Service, error) {
	o, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(instance.Cluster),
		Services: aws.StringSlice([]string{instance.Service}),
	})
	if err != nil {
		return nil, nil, err
	}

	if len(o.Services) == 0 {
		return nil, nil, fmt.Errorf("service not found with name %s", instance.Service)
	}

	if len(o.Services) > 1 {
		return nil, nil, fmt.Errorf("too many services returned for %s", instance.Service)
	}

	tdo, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(*o.Services[0].TaskDefinition),
	})
	if err != nil {
		return nil, nil, err
	}

	return tdo.TaskDefinition, o.Services[0], nil
}
