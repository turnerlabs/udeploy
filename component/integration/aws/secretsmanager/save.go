package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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
