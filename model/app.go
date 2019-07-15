package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Application ...
type Application struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	Name string     `json:"name" bson:"name"`
	Type string     `json:"type" bson:"type"`
	Repo Repository `json:"repo" bson:"repository"`

	Instances map[string]Instance `json:"instances" bson:"instances"`
}

// ToView ...
func (a Application) ToView(usr User) AppView {
	view := AppView{
		ID:        a.ID,
		Name:      a.Name,
		Type:      a.Type,
		Repo:      a.Repo,
		Instances: []InstanceView{},
	}

	for name, inst := range a.Instances {
		view.Instances = append(view.Instances, inst.ToView(name, usr.Apps[a.Name]))
	}

	return view
}

// GetInstances ...
func (a Application) GetInstances(filter []string) map[string]Instance {
	instances := map[string]Instance{}

	if len(filter) == 0 {
		return a.Instances
	}

	for _, ds := range filter {
		instances[ds] = a.Instances[ds]
	}

	return instances
}

// AppView ...
type AppView struct {
	ID primitive.ObjectID `json:"id,omitempty"`

	Name string     `json:"name"`
	Type string     `json:"type"`
	Repo Repository `json:"repo"`

	Instances []InstanceView `json:"instances"`
}

// ToBusiness ...
func (a AppView) ToBusiness() Application {
	app := Application{
		ID:        a.ID,
		Name:      a.Name,
		Type:      a.Type,
		Repo:      a.Repo,
		Instances: map[string]Instance{},
	}

	for _, inst := range a.Instances {
		app.Instances[inst.Name] = inst.ToBusiness()
	}

	return app
}

// Repository ...
type Repository struct {
	Org  string `json:"org" bson:"org"`
	Name string `json:"name" bson:"name"`

	AccessToken string `json:"accessToken" bson:"accessToken"`

	CommitConfig CommitConfig `json:"commitConfig" bson:"commitConfig"`
}

// CommitConfig ...
type CommitConfig struct {
	ExistingValue string `json:"existingValue" bson:"existingValue"`
	NewValue      string `json:"newValue" bson:"newValue"`

	Filter string `json:"filter" bson:"filter"`
	Limit  int    `json:"limit" bson:"limit"`
}
