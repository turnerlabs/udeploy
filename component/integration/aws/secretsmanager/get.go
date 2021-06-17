package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/turnerlabs/udeploy/component/app"
)

// Get ...
func Get(key string) (string, error) {
	return get(key, session.New())
}

// GetForInstance ...
func GetForInstance(key string, inst app.Instance) (string, error) {
	session := session.New()

	if len(inst.ConfigRole) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, inst.ConfigRole))
	}

	return get(key, session)

}

func get(key string, sess *session.Session) (string, error) {

	svc := secretsmanager.New(sess)

	o, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
				return "", nil
			}
		}

		return "", err
	}

	return *o.SecretString, nil
}
