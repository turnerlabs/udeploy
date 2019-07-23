package action

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// Complete ...
	Complete = "complete"

	// Pending ...
	Pending = "pending"

	// Error ...
	Error = "error"
)

// Action ...
type Action struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	DefinitionID string `json:"definitionId" bson:"definitionId"`

	Status string `json:"status" bson:"status"`
	Info   string `json:"info" bson:"info"`
	Type   string `json:"type" bson:"type"`

	Started time.Time `json:"started" bson:"started"`
	Stopped time.Time `json:"stopped" bson:"stopped"`
}

// Is ...
func (a Action) Is(status string) bool {
	return a.Status == status
}
