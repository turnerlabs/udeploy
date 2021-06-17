package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/turnerlabs/udeploy/component/app"
)

// Save ...
func Save(key, value, description, KmsKeyID string) error {

	svc := secretsmanager.New(session.New())

	_, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
				if _, err := svc.CreateSecret(&secretsmanager.CreateSecretInput{
					Name:         aws.String(key),
					SecretString: aws.String(value),
					Description:  aws.String(description),
					KmsKeyId:     aws.String(KmsKeyID),
				}); err != nil {
					return err
				}

				return nil
			}
		}

		return err
	}

	if _, err := svc.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(value),
		Description:  aws.String(description),
		KmsKeyId:     aws.String(KmsKeyID),
	}); err != nil {
		return err
	}

	return err
}

// UpdateForInstance ...
func UpdateForInstance(key, value string, inst app.Instance) error {
	session := session.New()

	if len(inst.ConfigRole) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, inst.ConfigRole))
	}

	svc := secretsmanager.New(session)

	_, err := svc.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(value),
	})

	return err
}
