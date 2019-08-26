package notice

import (
	"github.com/turnerlabs/udeploy/component/app"
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
	Error    bool `json:"error" bson:"error"`
	Pending  bool `json:"pending" bson:"pending"`
	Running  bool `json:"running" bson:"running"`
	Stopped  bool `json:"stopped" bson:"stopped"`
	Deployed bool `json:"deployed" bson:"deployed"`
}

// Matches ...
func (n Notice) Matches(instance string, t string, change app.Change) bool {

	if len(n.Instances) == 0 {
		return matchesEvent(n, t, change)
	}

	for _, i := range n.Instances {
		if i.Name == instance {
			return matchesEvent(n, t, change)
		}
	}

	return false
}

func matchesEvent(n Notice, t string, change app.Change) bool {

	switch t {
	case app.ChangeTypeVersion:
		if n.Events.Deployed {
			return true
		}
	case app.ChangeTypeStatus:
		switch change.After {
		case app.Pending:
			if n.Events.Pending {
				return true
			}
		case app.Running:
			if n.Events.Running {
				return true
			}
		case app.Stopped:
			if n.Events.Stopped {
				return true
			}
		}
	case app.ChangeTypeError:
		if n.Events.Error {
			return true
		}
	}

	return false
}
