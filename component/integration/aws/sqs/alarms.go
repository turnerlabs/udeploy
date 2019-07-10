package sqs

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/turnerlabs/udeploy/component/cfg"
	"go.mongodb.org/mongo-driver/mongo"
)

// MonitorAlarms ...
func MonitorAlarms(ctx mongo.SessionContext, fn process) error {

	svc := sqs.New(session.New())

	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(cfg.Get["SQS_ALARM_QUEUE"]),
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
			msg := alert{}

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
				json, err := msg.Message.MarshalJSON()
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

type alert struct {
	Message json.RawMessage `json:"Message"`
}

type alertMessage struct {
	Trigger alertTrigger `json:"Trigger"`

	NewStateReason string `json:"NewStateReason"`
	NewStateValue  string `json:"NewStateValue"`
	OldStateValue  string `json:"OldStateValue"`
}

type alertTrigger struct {
	Dimensions []alertMessageDimension `json:"Dimensions"`
}

type alertMessageDimension struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

func (a alert) toView() (MessageView, error) {
	msg := alertMessage{}

	s, _ := strconv.Unquote(string(a.Message))

	err := json.Unmarshal([]byte(s), &msg)
	if err != nil {
		return MessageView{}, err
	}

	name, alias := "", ""
	for _, d := range msg.Trigger.Dimensions {
		if d.Name == "Resource" {
			values := strings.Split(d.Value, ":")

			if len(values) == 2 {
				name = values[0]
				alias = values[1]
			}
		}
	}

	v := MessageView{
		ID: fmt.Sprintf("%s-%s", name, alias),
	}

	return v, nil
}
