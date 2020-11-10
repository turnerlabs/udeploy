package event

import (
	"fmt"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Populate ...
func Populate(instances map[string]app.Instance) (map[string]app.Instance, error) {

	populated := map[string]app.Instance{}

	for name, i := range instances {
		session := session.New()

		if len(i.Role) > 0 {
			session.Config.WithCredentials(stscreds.NewCredentials(session, i.Role))
		}

		evtSvc := cloudwatchevents.New(session)
		ecsSvc := ecs.New(session)

		i, state, err := populateInst(i, ecsSvc, evtSvc)
		if err != nil {
			state.SetError(err)
		}

		i.SetState(state)

		populated[name] = i
	}

	return populated, nil
}

func populateInst(i app.Instance, ecsSvc *ecs.ECS, evtSvc *cloudwatchevents.CloudWatchEvents) (app.Instance, app.State, error) {
	i.Task.Definition = app.NewDefinition(i.Task.Family)

	state := app.NewState()

	td, ruleOutput, target, err := getServiceInfo(i, ecsSvc, evtSvc)
	if err != nil {
		return i, state, err
	}

	i.Task.Definition, err = app.DefinitionFrom(td, i.Task.ImageTagEx)
	if err != nil {
		state.SetError(err)
	}

	runningTasks, err := task.List(i, ecsSvc, "RUNNING")
	if err != nil {
		return i, state, err
	}

	stoppedTasks, err := task.List(i, ecsSvc, "STOPPED")
	if err != nil {
		return i, state, err
	}

	isPending, isRunning, err := getStatus(i, td, runningTasks, stoppedTasks)

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

	state.Version = i.Task.Definition.Version.Full()

	i.Task.DesiredCount = *target.EcsParameters.TaskCount
	i.Task.CronEnabled = isCronEnabled(*ruleOutput.State)
	i.Task.CronExpression = fmt.Sprintf("0 %s", (*ruleOutput.ScheduleExpression)[5:len(*ruleOutput.ScheduleExpression)-1])

	i.Task.TasksInfo, err = task.GetTasksInfo(i, ecsSvc, append(runningTasks, stoppedTasks...))

	return i, state, err
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
				return false, false, app.InstanceError{Problem: fmt.Sprintf("container exited with code %d", *c.ExitCode)}
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
