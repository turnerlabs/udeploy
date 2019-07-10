package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"github.com/turnerlabs/udeploy/model"
)

// Deploy ...
func Deploy(source model.Instance, target model.Instance, revision int64, opts task.DeployOptions) error {

	svc := ecs.New(session.New())

	newOutput, err := task.Deploy(source, target, revision, source.Version(), opts)
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
