package sqs

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/turnerlabs/udeploy/component/cfg"
	"go.mongodb.org/mongo-driver/mongo"
)

// MonitorS3 ...
func MonitorS3(ctx mongo.SessionContext, fn process) error {

	svc := sqs.New(session.New())

	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(cfg.Get["SQS_S3_QUEUE"]),
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
			msg := s3msg{}

			if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
				log.Println(err)
				continue
			}

			if len(msg.Records) == 0 {
				continue
			}

			messagesToDelete = append(messagesToDelete, &sqs.DeleteMessageBatchRequestEntry{
				Id:            m.MessageId,
				ReceiptHandle: m.ReceiptHandle,
			})

			view, err := msg.toView()
			if err != nil {
				if err.Error() != unmonitoredMessageError {
					log.Println(msg)
					log.Println(err)
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

type s3msg struct {
	Records []msgRecord `json:"records"`
}

type msgRecord struct {
	S3 msgS3 `json:"s3"`
}

type msgS3 struct {
	Bucket msgBucket `json:"bucket"`
	Object msgObject `json:"object"`
}

type msgBucket struct {
	Name string `json:"name"`
}

type msgObject struct {
	Key string `json:"key"`
}

func (m s3msg) toView() (MessageView, error) {
	return MessageView{
		ID: fmt.Sprintf("%s-%s", m.Records[0].S3.Bucket.Name, m.Records[0].S3.Object.Key),
	}, nil
}
