package service

import (
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Scale ...
func Scale(instance app.Instance, desiredCount int64, restart bool) error {

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
