package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Deploy ...
func Deploy(source app.Instance, target app.Instance, revision int64, opts task.DeployOptions) error {

	session := session.New()

	config.Merge([]string{target.Role}, session)

	svc := ecs.New(session)

	newOutput, err := task.Deploy(source, target, revision, source.Task.Definition.Version.Version, opts)
	if err != nil {
		return err
	}

	_, err = svc.UpdateService(
		&ecs.UpdateServiceInput{
			Cluster:        aws.String(target.Cluster),
			Service:        aws.String(target.Service),
			TaskDefinition: newOutput.TaskDefinitionArn,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
