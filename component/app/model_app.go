package app

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/project"
	"github.com/turnerlabs/udeploy/component/user"
	"github.com/turnerlabs/udeploy/component/version"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AppTypeService       = "service"
	AppTypeScheduledTask = "scheduled-task"
	AppTypeLambda        = "lambda"
	AppTypeS3            = "s3"

	// FailedTaskExpiration determines the minutes an AWS ECS failed tasks should be considered.
	//
	// Currently the AWS ECS container restart throttle may wait a maximum of 15 minutes before
	// attempting a restart. 20 minutes has been chosen as the time to consider failed tasks.
	// This ensures a service does not appear to be healthy due to lack of attempted restarts
	// while reducing the time it takes to determine a service is healthy from 1 hour
	// down to 20 minutes.
	//
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-throttle-logic.html
	FailedTaskExpiration = 20

	Undetermined = "undetermined"
)

// Application ...
type Application struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	Name string     `json:"name" bson:"name"`
	Type string     `json:"type" bson:"type"`
	Repo Repository `json:"repo" bson:"repository"`

	ProjectID primitive.ObjectID `json:"projectId" bson:"projectId"`

	Instances map[string]Instance `json:"instances" bson:"instances"`
}

// ToView ...
func (a Application) ToView(usr user.User, project project.Project) AppView {
	view := AppView{
		ID:   a.ID,
		Name: a.Name,
		Type: a.Type,
		Repo: a.Repo,
		Project: ProjectView{
			ID:   a.ProjectID,
			Name: project.Name,
		},
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

// GetErrorInstances ...
func (a Application) GetErrorInstances() map[string]Instance {
	instances := map[string]Instance{}

	for k, i := range a.Instances {
		if i.CurrentState.Error != nil {
			instances[k] = i
		}
	}

	return instances
}

// Matches ...
func (a Application) Matches(filter Filter, p project.Project) bool {
	for _, t := range filter.Terms {
		if strings.Contains(strings.ToLower(a.Name), strings.ToLower(t)) ||
			(len(p.Name) > 0 && strings.Contains(strings.ToLower(p.Name), strings.ToLower(t))) {
			return true
		}
	}

	return len(filter.Terms) == 0
}

// Filter ...
type Filter struct {
	Terms []string `json:"terms"`
}

// AppView ...
type AppView struct {
	ID primitive.ObjectID `json:"id,omitempty"`

	Name string     `json:"name"`
	Type string     `json:"type"`
	Repo Repository `json:"repo"`

	Project ProjectView `json:"project"`

	Instances []InstanceView `json:"instances"`
}

// ProjectView ...
type ProjectView struct {
	ID primitive.ObjectID `json:"id,omitempty"`

	Name string `json:"name"`
}

// ToBusiness ...
func (a AppView) ToBusiness() Application {
	app := Application{
		ID:        a.ID,
		Name:      a.Name,
		Type:      a.Type,
		Repo:      a.Repo,
		ProjectID: a.Project.ID,
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

// Definition ...
type Definition struct {
	ID string `json:"id"`

	Version  string `json:"version"`
	Build    string `json:"build"`
	Revision int64  `json:"revision"`

	Description string `json:"description"`

	Environment map[string]string `json:"environment"`
	Secrets     map[string]string `json:"secrets"`

	Registry bool `json:"registry"`
}

// FormatVersion ...
func (d Definition) FormatVersion() string {

	if d.Version == "" {
		return Undetermined
	}

	if len(d.Build) > 0 {
		return fmt.Sprintf("%s.%s", d.Version, d.Build)
	}

	return d.Version
}

// NewDefinition ...
func NewDefinition(id string) Definition {
	return Definition{
		ID: id,

		Environment: map[string]string{},
		Secrets:     map[string]string{},
	}
}

// DefinitionFrom ...
func DefinitionFrom(td *ecs.TaskDefinition, imageTagRegEx string) Definition {

	version, build := version.Extract(*td.ContainerDefinitions[0].Image, imageTagRegEx)

	def := Definition{
		ID: (*td.TaskDefinitionArn)[0:strings.LastIndex(*td.TaskDefinitionArn, ":")],

		Version:  version,
		Build:    build,
		Revision: *td.Revision,

		Description: *td.ContainerDefinitions[0].Image,

		Environment: map[string]string{},
		Secrets:     map[string]string{},
	}

	for _, e := range td.ContainerDefinitions[0].Environment {
		value := *e.Value
		def.Environment[*e.Name] = value
	}

	for _, e := range td.ContainerDefinitions[0].Secrets {
		value := *e.ValueFrom
		def.Secrets[*e.Name] = value
	}

	return def
}
