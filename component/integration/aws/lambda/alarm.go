package lambda

import (
	"fmt"

	"github.com/turnerlabs/udeploy/component/cfg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// DeleteAlarm ...
func DeleteAlarm(functionName, role string) error {
	session := session.New()

	if len(role) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, role))
	}

	svc := cloudwatch.New(session)

	_, err := svc.DeleteAlarms(&cloudwatch.DeleteAlarmsInput{
		AlarmNames: aws.StringSlice([]string{buildAlarmName(functionName)}),
	})

	return err
}

// UpsertAlarm ...
func UpsertAlarm(functionName, functionAlias, role, snsTopicArn string) error {
	session := session.New()

	if len(role) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, role))
	}

	svc := cloudwatch.New(session)

	_, err := svc.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(buildAlarmName(functionName)),
		AlarmDescription:   aws.String(fmt.Sprintf("monitor %s errors", functionName)),
		ComparisonOperator: aws.String("GreaterThanThreshold"),
		EvaluationPeriods:  aws.Int64(1),
		MetricName:         aws.String("Errors"),
		Namespace:          aws.String("AWS/Lambda"),
		Period:             aws.Int64(60),
		Statistic:          aws.String("Maximum"),
		Threshold:          aws.Float64(0),
		TreatMissingData:   aws.String("ignore"),
		Dimensions: []*cloudwatch.Dimension{&cloudwatch.Dimension{
			Name:  aws.String("FunctionName"),
			Value: aws.String(functionName),
		}, &cloudwatch.Dimension{
			Name:  aws.String("Resource"),
			Value: aws.String(fmt.Sprintf("%s:%s", functionName, functionAlias)),
		},
		},

		AlarmActions:            aws.StringSlice([]string{snsTopicArn}),
		OKActions:               aws.StringSlice([]string{snsTopicArn}),
		InsufficientDataActions: aws.StringSlice([]string{snsTopicArn}),
	})

	return err
}

func buildAlarmName(functionName string) string {
	return fmt.Sprintf("%s-%s-lambda-errors", functionName, cfg.Get["ENV"])
}
