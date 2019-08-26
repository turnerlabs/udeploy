package event

import (
	"fmt"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Populate ...
func Populate(instances map[string]app.Instance) (map[string]app.Instance, error) {
	sesssion := session.New()

	evtSvc := cloudwatchevents.New(sesssion)
	ecsSvc := ecs.New(sesssion)

	instanceChan := make(chan chanModel, len(instances))

	for name, instance := range instances {

		go func(innerName string, innerInstance app.Instance, innerEcs *ecs.ECS, innerEvt *cloudwatchevents.CloudWatchEvents) {
			state := app.NewState()

			td, ruleOutput, target, err := getServiceInfo(innerInstance, innerEcs, innerEvt)
			if err != nil {
				state.SetError(err)
			} else {
				innerInstance.Task.Definition = app.DefinitionFrom(td, innerInstance.Task.ImageTagEx)

				runningTasks, err := task.List(innerInstance, innerEcs, "RUNNING")
				if err != nil {
					state.SetError(err)
				}

				stoppedTasks, err := task.List(innerInstance, innerEcs, "STOPPED")
				if err != nil {
					state.SetError(err)
				}

				isPending, isRunning, err := getStatus(innerInstance, td, runningTasks, stoppedTasks)

				if err != nil {
					state.SetError(err)
				}

				if isPending {
					state.SetPending()
				} else if isRunning {
					state.SetRunning()
				} else {
					state.SetStopped()
				}

				state.Version = innerInstance.FormatVersion()

				innerInstance.Task.DesiredCount = *target.EcsParameters.TaskCount
				innerInstance.Task.CronEnabled = isCronEnabled(*ruleOutput.State)
				innerInstance.Task.CronExpression = fmt.Sprintf("0 %s", (*ruleOutput.ScheduleExpression)[5:len(*ruleOutput.ScheduleExpression)-1])

				innerInstance.Task.TasksInfo, err = task.GetTasksInfo(innerInstance, innerEcs, append(runningTasks, stoppedTasks...))
				if err != nil {
					state.SetError(err)
				}
			}

			innerInstance.SetState(state)

			instanceChan <- chanModel{
				name:     innerName,
				instance: innerInstance,
			}

		}(name, instance, ecsSvc, evtSvc)
	}

	for respCount := 1; respCount <= len(instances); respCount++ {
		i := <-instanceChan

		instances[i.name] = i.instance

		if respCount == len(instances) {
			close(instanceChan)
		}
	}

	return instances, nil
}

func isCronEnabled(state string) bool {
	return state == "ENABLED"
}

type chanModel struct {
	name     string
	instance app.Instance
}

func getStatus(instance app.Instance, td *ecs.TaskDefinition, runningTasks, stoppedTasks []*ecs.Task) (isPending bool, isRunning bool, err error) {
	tasks := append(runningTasks, stoppedTasks...)

	lastTask := &ecs.Task{}
	for i, t := range tasks {
		if i == 0 || (*t.CreatedAt).After(*lastTask.CreatedAt) {
			lastTask = t
		}
	}

	if lastTask != nil && lastTask.LastStatus != nil && *lastTask.LastStatus != "RUNNING" {
		for _, c := range lastTask.Containers {
			if c.ExitCode != nil && *c.ExitCode != 0 {
				return false, false, fmt.Errorf("container exited with code %d", *c.ExitCode)
			}
		}
	}

	if len(runningTasks) == 0 {
		return false, false, nil
	}

	for _, task := range runningTasks {
		if *task.DesiredStatus != "RUNNING" || *task.LastStatus != "RUNNING" {
			return true, false, nil
		}
	}

	return false, true, nil
}

func getServiceInfo(instance app.Instance, ecsSvc *ecs.ECS, evtSvc *cloudwatchevents.CloudWatchEvents) (*ecs.TaskDefinition, *cloudwatchevents.DescribeRuleOutput, *cloudwatchevents.Target, error) {
	ruleInput := &cloudwatchevents.DescribeRuleInput{
		Name: &instance.EventRule,
	}
	ruleOutput, err := evtSvc.DescribeRule(ruleInput)
	if err != nil {
		return nil, nil, nil, err
	}

	targetInput := &cloudwatchevents.ListTargetsByRuleInput{
		Rule: &instance.EventRule,
	}
	targetOutput, err := evtSvc.ListTargetsByRule(targetInput)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(targetOutput.Targets) == 0 {
		return nil, nil, nil, fmt.Errorf("event rule target not found")
	}

	target := targetOutput.Targets[0]

	tdo, err := ecsSvc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: target.EcsParameters.TaskDefinitionArn,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return tdo.TaskDefinition, ruleOutput, target, nil
}
