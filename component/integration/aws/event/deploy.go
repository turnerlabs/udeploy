package event

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"github.com/turnerlabs/udeploy/model"
)

func Deploy(source model.Instance, target model.Instance, revision int64, opts task.DeployOptions) error {
	svc := cloudwatchevents.New(session.New())

	newOutput, err := task.Deploy(source, target, revision, source.Version(), opts)
	if err != nil {
		return err
	}

	err = updateTargetRevision(target, svc, newOutput.TaskDefinitionArn)
	if err != nil {
		return err
	}

	return nil
}

func updateTargetRevision(instance model.Instance, svc *cloudwatchevents.CloudWatchEvents, taskArn *string) error {
	input := &cloudwatchevents.ListTargetsByRuleInput{
		Rule: &instance.EventRule,
	}
	resp, err := svc.ListTargetsByRule(input)
	if err != nil {
		return err
	}
	if len(resp.Targets) == 0 {
		return fmt.Errorf("event target not found")
	}

	if len(resp.Targets) > 1 {
		return fmt.Errorf("too many targets")
	}

	resp.Targets[0].EcsParameters.TaskDefinitionArn = taskArn

	putInput := &cloudwatchevents.PutTargetsInput{
		Rule:    &instance.EventRule,
		Targets: resp.Targets,
	}
	putResp, err := svc.PutTargets(putInput)
	if err != nil {
		return err
	}
	if putResp != nil {
		if len(putResp.FailedEntries) > 0 {
			for _, entry := range putResp.FailedEntries {
				return fmt.Errorf("TargetId: %s; ErrorCode: %s; ErrorMessage: %s", *entry.TargetId, *entry.ErrorCode, *entry.ErrorMessage)
			}
		}
	}
	return nil
}
