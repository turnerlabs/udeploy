package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/model"
)

// Scale ...
func Scale(instance model.Instance, desiredCount int64, restart bool) error {

	svc := ecs.New(session.New())

	_, err := svc.UpdateService(
		&ecs.UpdateServiceInput{
			Cluster:            aws.String(instance.Cluster),
			Service:            aws.String(instance.Service),
			DesiredCount:       aws.Int64(desiredCount),
			ForceNewDeployment: aws.Bool(restart),
		},
	)

	if err != nil {
		return err
	}

	return nil
}
