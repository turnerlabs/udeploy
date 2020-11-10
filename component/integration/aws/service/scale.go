package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/app"
)

// Scale ...
func Scale(instance app.Instance, desiredCount int64, restart bool) error {

	session := session.New()

	if len(instance.Role) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, instance.Role))
	}

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
