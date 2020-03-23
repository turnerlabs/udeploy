package project

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project ...
type Project struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

// FindByID ...
func FindByID(ID primitive.ObjectID, projects []Project) (Project, bool) {
	for _, p := range projects {
		if ID == p.ID {
			return p, true
		}
	}

	return Project{}, false
}
