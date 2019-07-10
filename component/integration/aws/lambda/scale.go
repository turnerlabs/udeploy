package lambda

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turnerlabs/udeploy/model"
)

// Scale ...
func Scale(ctx context.Context, instance model.Instance, desiredCount int64) error {

	if desiredCount == 0 {
		return nil
	}

	svc := lambda.New(session.New())

	ao, err := svc.GetAliasWithContext(ctx, &lambda.GetAliasInput{
		Name:         aws.String(instance.FunctionAlias),
		FunctionName: aws.String(instance.FunctionName),
	})
	if err != nil {
		return err
	}

	for x := int64(1); x <= desiredCount; x++ {
		_, err := svc.InvokeWithContext(ctx,
			&lambda.InvokeInput{
				FunctionName:   ao.AliasArn,
				InvocationType: aws.String(lambda.InvocationTypeEvent),
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}
