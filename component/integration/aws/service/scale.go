package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
)

// Scale ...
func Scale(instance app.Instance, desiredCount int64, restart bool) error {

	session := session.New()

	config.Merge([]string{instance.Role}, session)

	svc := ecs.New(session)

	_, err := svc.UpdateService(
		&ecs.UpdateServiceInput{
			Cluster:            aws.String(instance.Cluster),
			Service:            aws.String(instance.Service),
			DesiredCount:       aws.Int64(desiredCount),
			ForceNewDeployment: aws.Bool(restart),
		},
	)

	return err
}
