package sqs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/turnerlabs/udeploy/component/cfg"
)

// MonitorChanges ...
func MonitorChanges(ctx mongo.SessionContext, fn process) error {

	svc := sqs.New(session.New())

	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(cfg.Get["SQS_CHANGE_QUEUE"]),
	})
	if err != nil {
		return err
	}

	for {
		o, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            url.QueueUrl,
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(timeout),
		})
		if err != nil {
			return err
		}

		messagesToDelete := []*sqs.DeleteMessageBatchRequestEntry{}
		for _, m := range o.Messages {
			msg := message{}

			if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
				log.Println(err)
				continue
			}

			messagesToDelete = append(messagesToDelete, &sqs.DeleteMessageBatchRequestEntry{
				Id:            m.MessageId,
				ReceiptHandle: m.ReceiptHandle,
			})

			view, err := msg.toView()
			if err != nil {
				if err.Error() == unmonitoredMessageError {
					continue
				}

				log.Println(err)
				json, err := msg.Detail.MarshalJSON()
				if err == nil {
					log.Println(string(json))
				}

				continue
			}

			if err := fn(ctx, view); err != nil {
				log.Println(err)
				continue
			}
		}

		if len(messagesToDelete) > 0 {
			if _, err := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
				QueueUrl: url.QueueUrl,
				Entries:  messagesToDelete,
			}); err != nil {
				log.Println(err)
			}
		}
	}
}

type message struct {
	Type   string          `json:"detail-type"`
	Source string          `json:"source"`
	Detail json.RawMessage `json:"detail"`
}

func (m message) toView() (MessageView, error) {

	switch m.Type {
	case "AWS API Call via CloudTrail":

		switch m.Source {
		case "aws.lambda":
			detail := lambdaMessageDetail{}

			if err := json.Unmarshal(m.Detail, &detail); err != nil {
				return MessageView{}, err
			}

			if detail.ErrorCode == "ClientException" {
				return MessageView{}, errors.New(unmonitoredMessageError)
			}

			if !strings.Contains(detail.EventName, "UpdateAlias") {
				return MessageView{}, errors.New(unmonitoredMessageError)
			}

			return MessageView{
				ID: detail.taskDefinition(),
			}, nil
		default:
			detail := logMessageDetail{}

			if err := json.Unmarshal(m.Detail, &detail); err != nil {
				return MessageView{}, err
			}

			switch detail.ErrorCode {
			case "ClientException", "AccessDenied":
				return MessageView{}, errors.New(unmonitoredMessageError)
			}

			id, err := detail.taskDefinition()
			if err != nil {
				return MessageView{}, err
			}

			return MessageView{
				ID: id,
			}, nil
		}
	case "ECS Task State Change":
		detail := eventMessageDetail{}

		if err := json.Unmarshal(m.Detail, &detail); err != nil {
			return MessageView{}, err
		}

		return MessageView{
			ID: detail.taskDefinition(),
		}, nil
	}

	return MessageView{}, nil
}

type lambdaMessageDetail struct {
	EventName         string            `json:"eventName"`
	ErrorCode         string            `json:"errorCode"`
	RequestParameters requestParameters `json:"requestParameters"`
}

type requestParameters struct {
	FunctionName string `json:"functionName"`
	AliasName    string `json:"name"`
}

func (d lambdaMessageDetail) taskDefinition() string {
	return fmt.Sprintf("%s-%s", d.RequestParameters.FunctionName, d.RequestParameters.AliasName)
}

type logMessageDetail struct {
	ErrorCode   string      `json:"errorCode"`
	LogResponse logResponse `json:"responseElements"`
}

func (d logMessageDetail) taskDefinition() (string, error) {
	if len(d.LogResponse.LogTaskDefinition.TaskDefinitionArn) == 0 {
		return "", errors.New("task definition not found")
	}

	return d.LogResponse.LogTaskDefinition.TaskDefinitionArn[strings.Index(d.LogResponse.LogTaskDefinition.TaskDefinitionArn, "/")+1 : strings.LastIndex(d.LogResponse.LogTaskDefinition.TaskDefinitionArn, ":")], nil
}

type logResponse struct {
	LogTaskDefinition logTaskDefinition `json:"taskDefinition"`
}

type logTaskDefinition struct {
	TaskDefinitionArn string `json:"taskDefinitionArn"`
}

type eventMessageDetail struct {
	TaskDefinitionArn string `json:"taskDefinitionArn"`
	LastStatus        string `json:"lastStatus"`
	DesiredStatus     string `json:"desiredStatus"`
}

func (d eventMessageDetail) taskDefinition() string {
	return d.TaskDefinitionArn
}
