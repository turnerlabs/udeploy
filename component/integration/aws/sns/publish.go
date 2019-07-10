package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Publish ...
func Publish(subject, msg, arn string) error {
	svc := sns.New(session.New())

	_, err := svc.Publish(&sns.PublishInput{
		Subject:  aws.String(subject),
		TopicArn: aws.String(arn),
		Message:  aws.String(msg),
	})
	if err != nil {
		return err
	}

	return nil
}
