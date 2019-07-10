package sqs

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const timeout int64 = 20

var unmonitoredMessageError = "unmonitored message"

type process func(mongo.SessionContext, MessageView) error

// MessageView ...
type MessageView struct {
	ID string
}
