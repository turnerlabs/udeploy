package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notice ...
type Notice struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Enabled   bool               `json:"enabled" bson:"enabled"`
	SNSArn    string             `json:"snsArn" bson:"snsArn"`
	Apps      []NoticeOption     `json:"apps" bson:"apps"`
	Instances []NoticeOption     `json:"instances" bson:"instances"`
	Events    NoticeEvents       `json:"events" bson:"events"`
}

// NoticeOption ...
type NoticeOption struct {
	Name string `json:"name" bson:"name"`
}

// NoticeEvents ...
type NoticeEvents struct {
	Error     bool `json:"error" bson:"error"`
	Starting  bool `json:"starting" bson:"starting"`
	Pending   bool `json:"pending" bson:"pending"`
	Running   bool `json:"running" bson:"running"`
	Stopped   bool `json:"stopped" bson:"stopped"`
	Deployed  bool `json:"deployed" bson:"deployed"`
	Deploying bool `json:"deploying" bson:"deploying"`
}

// Matches ...
func (n Notice) Matches(instance string, inst Instance) bool {

	switch inst.String() {
	case "error":
		if !n.Events.Error {
			return false
		}
	case "starting":
		if !n.Events.Starting {
			return false
		}
	case "pending":
		if !n.Events.Pending {
			return false
		}
	case "running":
		if !n.Events.Running {
			return false
		}
	case "stopped":
		if !n.Events.Stopped {
			return false
		}
	case "deployed":
		if !n.Events.Deployed {
			return false
		}
	case "deploying":
		if !n.Events.Deploying {
			return false
		}
	}

	if len(n.Instances) == 0 {
		return true
	}

	for _, i := range n.Instances {
		if instance == i.Name {
			return true
		}
	}

	return false
}
