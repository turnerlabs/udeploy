package event

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

func Deploy(source app.Instance, target app.Instance, revision int64, opts task.DeployOptions) error {

	newOutput, err := task.Deploy(source, target, revision, source.Version(), opts)
	if err != nil {
		return err
	}

	session := session.New()

	config.Merge([]string{source.Role, target.Role}, session)

	svc := cloudwatchevents.New(session)

	return updateTargetRevision(target, svc, newOutput.TaskDefinitionArn)
}

func updateTargetRevision(instance app.Instance, svc *cloudwatchevents.CloudWatchEvents, taskArn *string) error {
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
