package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Delete ...
func Delete(key string) error {
	svc := secretsmanager.New(session.New())

	_, err := svc.DeleteSecret(&secretsmanager.DeleteSecretInput{
		ForceDeleteWithoutRecovery: aws.Bool(true),
		SecretId:                   aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
				return nil
			}
		}

		return err
	}

	return nil
}
