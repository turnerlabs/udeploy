package event

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"github.com/turnerlabs/udeploy/model"
)

// Populate ...
func Populate(instances map[string]model.Instance, details bool) (map[string]model.Instance, error) {
	sesssion := session.New()

	evtSvc := cloudwatchevents.New(sesssion)
	ecsSvc := ecs.New(sesssion)

	instanceChan := make(chan chanModel, len(instances))

	for name, instance := range instances {

		go func(innerName string, innerInstance model.Instance, innerEcs *ecs.ECS, innerEvt *cloudwatchevents.CloudWatchEvents) {
			state := model.State{}

			td, ruleOutput, target, err := getServiceInfo(innerInstance, innerEcs, innerEvt)
			if err != nil {
				state.Error = err
			} else {
				innerInstance.Task.Definition = model.DefinitionFrom(td, innerInstance.Task.ImageTagEx)

				state.IsPending, state.IsRunning, err = getStatus(innerInstance, td, innerEcs)
				if err != nil {
					state.Error = err
				}
				state.DesiredCount = *target.EcsParameters.TaskCount
				state.Version = innerInstance.FormatVersion()

				innerInstance.Task.CronEnabled = isCronEnabled(*ruleOutput.State)
				innerInstance.Task.CronExpression = fmt.Sprintf("0 %s", (*ruleOutput.ScheduleExpression)[5:len(*ruleOutput.ScheduleExpression)-1])

				if details {
					innerInstance.Task.TasksInfo, err = task.GetTasksInfo(innerInstance, innerEcs)
					if err != nil {
						state.Error = err
					}
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
	instance model.Instance
}

func getStatus(instance model.Instance, td *ecs.TaskDefinition, svc *ecs.ECS) (isPending bool, isRunning bool, err error) {
	tasks, err := task.List(instance, svc, aws.String("RUNNING"))
	if err != nil {
		return false, false, err
	}
	if len(tasks) == 0 {
		return false, false, nil
	}
	for _, task := range tasks {
		if *task.DesiredStatus != "RUNNING" || *task.LastStatus != "RUNNING" {
			return true, false, nil
		}
	}
	return false, true, nil
}

func getServiceInfo(instance model.Instance, ecsSvc *ecs.ECS, evtSvc *cloudwatchevents.CloudWatchEvents) (*ecs.TaskDefinition, *cloudwatchevents.DescribeRuleOutput, *cloudwatchevents.Target, error) {
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
