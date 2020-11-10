package event

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

func Deploy(source app.Instance, target app.Instance, revision int64, opts task.DeployOptions) error {

	newOutput, err := task.Deploy(source, target, revision, source.Task.Definition.Version.Version, opts)
	if err != nil {
		return err
	}

	return updateTargetRevision(target, newOutput.TaskDefinitionArn)
}

func updateTargetRevision(instance app.Instance, taskArn *string) error {
	session := session.New()

	if len(instance.Role) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, instance.Role))
	}

	svc := cloudwatchevents.New(session)

	resp, err := svc.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
		Rule: &instance.EventRule,
	})
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

	putResp, err := svc.PutTargets(&cloudwatchevents.PutTargetsInput{
		Rule:    &instance.EventRule,
		Targets: resp.Targets,
	})
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
