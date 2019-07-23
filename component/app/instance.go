package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turnerlabs/udeploy/component/user"
)

// InstanceView ...
type InstanceView struct {
	// Database Fields
	Name          string `json:"name"`
	Order         int    `json:"order"`
	Cluster       string `json:"cluster,omitempty"`
	Service       string `json:"service,omitempty"`
	EventRule     string `json:"eventRule,omitempty"`
	FunctionName  string `json:"functionName,omitempty"`
	FunctionAlias string `json:"functionAlias,omitempty"`

	S3Bucket         string `json:"s3Bucket,omitempty"`
	S3ConfigKey      string `json:"s3ConfigKey,omitempty"`
	S3Prefix         string `json:"s3Prefix,omitempty"`
	S3RegistryBucket string `json:"s3RegistryBucket,omitempty"`
	S3RegistryPrefix string `json:"s3RegistryPrefix,omitempty"`

	Repository    string   `json:"repository,omitempty"`
	DeployCode    string   `json:"deployCode"`
	AutoPropagate bool     `json:"autoPropagate"`
	Task          TaskView `json:"task"`

	// Calculated Fields
	FormattedVersion string `json:"formattedVersion"`
	Version          string `json:"version"`
	Revision         int64  `json:"revision"`

	Containers []ContainerView `json:"containers"`
	Deployment DeploymentView  `json:"deployment"`

	Error        string `json:"error"`
	IsRunning    bool   `json:"isRunning"`
	DesiredCount int64  `json:"desiredCount"`
	Edited       bool   `json:"edited"`

	CronExpression string `json:"cronExpression"`
	CronEnabled    bool   `json:"cronEnabled"`

	Claims map[string]bool   `json:"claims"`
	Tokens map[string]string `json:"tokens"`

	Links []Link `json:"links"`
}

// Link ...
type Link struct {
	Name        string `json:"name" bson:"name"`
	URL         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"decription"`
	Generated   bool   `json:"generated" bson:"-"`
}

// DeploymentView ...
type DeploymentView struct {
	IsPending bool `json:"isPending"`
}

// ContainerView ...
type ContainerView struct {
	Image       string            `json:"image"`
	Environment map[string]string `json:"environment"`
	Secrets     map[string]string `json:"secrets"`
}

// Instance ...
type Instance struct {
	// Database Fields
	Order   int    `json:"order" bson:"order"`
	Cluster string `json:"cluster,omitempty" bson:"cluster"`
	Service string `json:"service,omitempty" bson:"service"`

	EventRule string `json:"eventRule,omitempty" bson:"eventRule"`

	FunctionName  string `json:"functionName,omitempty" bson:"functionName"`
	FunctionAlias string `json:"functionAlias,omitempty" bson:"functionAlias"`

	S3Bucket         string `json:"s3Bucket,omitempty" bson:"s3Bucket"`
	S3ConfigKey      string `json:"s3ConfigKey,omitempty" bson:"s3ConfigKey"`
	S3Prefix         string `json:"s3Prefix,omitempty" bson:"s3Prefix"`
	S3RegistryBucket string `json:"s3RegistryBucket,omitempty" bson:"s3RegistryBucket"`
	S3RegistryPrefix string `json:"s3RegistryPrefix,omitempty" bson:"s3RegistryPrefix"`

	Repository    string `json:"repository,omitempty" bson:"repository"`
	DeployCode    string `json:"deployCode" bson:"deployCode"`
	AutoPropagate bool   `json:"autoPropagate" bson:"autoPropagate"`
	Task          Task   `json:"task" bson:"taskDefinition"`
	Links         []Link `json:"links" bson:"links"`

	// Calculated Fields
	CurrentState  State `json:"-" bson:"-"`
	previousState State
}

// S3FullConfigKey ...
func (i *Instance) S3FullConfigKey() string {
	if len(i.S3Prefix) > 0 {
		return fmt.Sprintf("%s/%s", i.S3Prefix, i.S3ConfigKey)
	}

	return i.S3ConfigKey
}

// SetState ...
func (i *Instance) SetState(s State) {
	i.previousState = s
	i.CurrentState = s
}

// RecordState ...
func (i *Instance) RecordState(s State) {
	i.previousState = s
}

// Changed ...
func (i *Instance) Changed() (bool, map[string]Change) {
	return i.CurrentState.ChangedFrom(i.previousState)
}

// Status ...
func (i *Instance) String() string {
	if !i.CurrentState.IsPending {
		if i.CurrentState.Version != i.previousState.Version {
			return "deployed"
		}
	}

	if i.CurrentState.Error != nil {
		return "error"
	}

	if i.CurrentState.IsPending {
		if i.CurrentState.Version != i.previousState.Version {
			return "deploying"
		}

		if i.CurrentState.IsRunning {
			return "pending"
		}

		if i.CurrentState.DesiredCount > 0 {
			return "starting"
		}
	}

	if i.CurrentState.IsRunning {
		return "running"
	}

	return "stopped"
}

// State ...
type State struct {
	Version      string `json:"-" bson:"-"`
	IsPending    bool   `json:"-" bson:"-"`
	IsRunning    bool   `json:"-" bson:"-"`
	DesiredCount int64  `json:"-" bson:"-"`
	Error        error  `json:"-" bson:"-"`
}

// Change ...
type Change struct {
	Before string
	After  string
}

func (c Change) String() string {
	return fmt.Sprintf("%s => %s", c.Before, c.After)
}

// ChangedFrom ...
func (s State) ChangedFrom(prev State) (bool, map[string]Change) {
	changes := map[string]Change{}

	if s.Error != nil && prev.Error == nil {
		changes["ERROR"] = Change{
			Before: fmt.Sprintf("%s", prev.Error),
			After:  fmt.Sprintf("%s", s.Error),
		}
	}

	if s.IsPending != prev.IsPending {
		changes["PENDING"] = Change{
			Before: strconv.FormatBool(prev.IsPending),
			After:  strconv.FormatBool(s.IsPending),
		}
	}

	if s.IsRunning != prev.IsRunning {
		changes["RUNNING"] = Change{
			Before: strconv.FormatBool(prev.IsRunning),
			After:  strconv.FormatBool(s.IsRunning),
		}
	}

	if s.DesiredCount != prev.DesiredCount {
		changes["COUNT"] = Change{
			Before: strconv.FormatInt(prev.DesiredCount, 10),
			After:  strconv.FormatInt(s.DesiredCount, 10),
		}
	}

	if s.Version != prev.Version {
		changes["VERSION"] = Change{
			Before: prev.Version,
			After:  s.Version,
		}
	}

	return len(changes) > 0, changes
}

// TaskView ...
type TaskView struct {
	// Database Fields
	Family       string `json:"family,omitempty"`
	Registry     string `json:"registry,omitempty"`
	ImageTagEx   string `json:"imageTagEx,omitempty"`
	CloneEnvVars string `json:"cloneEnvVars,omitempty"`
	Revisions    int    `json:"revisions,omitempty"`

	// Calculated Fields
	CronExpression string     `json:"cronExpression"`
	CronEnabled    bool       `json:"cronEnabled"`
	TasksInfo      []TaskInfo `json:"tasksInfo"`
}

// Task ...
type Task struct {
	// Database Fields
	Family       string   `json:"family,omitempty" bson:"family"`
	Registry     string   `json:"registry,omitempty" bson:"registry"`
	ImageTagEx   string   `json:"imageTagEx,omitempty" bson:"imageTagEx"`
	CloneEnvVars []string `json:"cloneEnvVars,omitempty" bson:"cloneEnvVars,omitempty"`
	Revisions    int      `json:"revisions,omitempty" bson:"revisions"`

	// Calculated Fields
	Definition     Definition `json:"-" bson:"-"`
	CronExpression string     `json:"-" bson:"-"`
	CronEnabled    bool       `json:"-" bson:"-"`
	TasksInfo      []TaskInfo `json:"-" bson:"-"`
}

type TaskInfo struct {
	TaskID         string    `json:"taskID"`
	LastStatus     string    `json:"lastStatus"`
	LastStatusTime time.Time `json:"lastStatusTime"`
	Version        string    `json:"version"`
	LogLink        string    `json:"logLink"`
	Reason         string    `json:"reason"`
}

// ToBusiness ...
func (v InstanceView) ToBusiness() Instance {

	i := Instance{
		Cluster:          v.Cluster,
		Service:          v.Service,
		EventRule:        v.EventRule,
		Repository:       v.Repository,
		FunctionName:     v.FunctionName,
		FunctionAlias:    v.FunctionAlias,
		S3Bucket:         v.S3Bucket,
		S3Prefix:         v.S3Prefix,
		S3ConfigKey:      v.S3ConfigKey,
		S3RegistryBucket: v.S3RegistryBucket,
		S3RegistryPrefix: v.S3RegistryPrefix,
		DeployCode:       v.DeployCode,
		AutoPropagate:    v.AutoPropagate,
		Order:            v.Order,
		Task: Task{
			ImageTagEx:   v.Task.ImageTagEx,
			CronEnabled:  v.Task.CronEnabled,
			Family:       v.Task.Family,
			Registry:     v.Task.Registry,
			Revisions:    v.Task.Revisions,
			TasksInfo:    v.Task.TasksInfo,
			CloneEnvVars: strings.Split(v.Task.CloneEnvVars, "\n"),
		},
	}

	for _, l := range v.Links {
		if !l.Generated {
			i.Links = append(i.Links, l)
		}
	}

	return i
}

// ToView ...
func (i Instance) ToView(name string, appClaim user.AppClaim) InstanceView {

	v := InstanceView{
		Name:             name,
		FormattedVersion: i.FormatVersion(),
		Version:          i.Version(),
		Containers:       []ContainerView{},
		IsRunning:        i.CurrentState.IsRunning,
		DesiredCount:     i.CurrentState.DesiredCount,
		CronExpression:   i.Task.CronExpression,
		CronEnabled:      i.Task.CronEnabled,
		Claims:           map[string]bool{},
		Cluster:          i.Cluster,
		Service:          i.Service,
		EventRule:        i.EventRule,
		FunctionName:     i.FunctionName,
		FunctionAlias:    i.FunctionAlias,
		S3Bucket:         i.S3Bucket,
		S3Prefix:         i.S3Prefix,
		S3ConfigKey:      i.S3ConfigKey,
		S3RegistryBucket: i.S3RegistryBucket,
		S3RegistryPrefix: i.S3RegistryPrefix,
		Repository:       i.Repository,
		DeployCode:       i.DeployCode,
		AutoPropagate:    i.AutoPropagate,
		Order:            i.Order,
		Links:            i.Links,
		Task: TaskView{
			ImageTagEx:   i.Task.ImageTagEx,
			CronEnabled:  i.Task.CronEnabled,
			Family:       i.Task.Family,
			Registry:     i.Task.Registry,
			Revisions:    i.Task.Revisions,
			TasksInfo:    i.Task.TasksInfo,
			CloneEnvVars: strings.Join(i.Task.CloneEnvVars, "\n"),
		},
		Deployment: DeploymentView{
			IsPending: i.CurrentState.IsPending,
		},
	}

	if v.Tokens == nil {
		v.Tokens = map[string]string{}
	}

	v.Tokens["VERSION"] = i.Version()

	if v.Links == nil {
		v.Links = []Link{}
	}

	for _, c := range appClaim.Claims[name] {
		v.Claims[c] = true
	}

	if i.CurrentState.Error != nil {
		v.Error = i.CurrentState.Error.Error()
	}

	v.Revision = i.Task.Definition.Revision

	v.Containers = append(v.Containers, ContainerView{
		Image:       i.Task.Definition.Description,
		Environment: i.Task.Definition.Environment,
		Secrets:     i.Task.Definition.Secrets,
	})

	return v
}

// Repo ...
func (i Instance) Repo() string {
	if index := strings.Index(i.Repository, ":"); index > -1 {
		return i.Repository[0:index]
	}

	return i.Repository
}

// RepoVersion ...
func (i Instance) RepoVersion() (string, string) {
	if index := strings.Index(i.Repository, ":"); index > -1 {
		version := i.Repository[index+1 : len(i.Repository)]
		parts := strings.Split(version, ".")

		if len(parts) > 1 {
			return parts[0], parts[1]
		}

		return parts[0], ""
	}

	return "", ""
}

// FormatVersion ...
func (i Instance) FormatVersion() string {

	if i.Task.Definition.Version == "" {
		return "undetermined"
	}

	if len(i.Task.Definition.Build) > 0 {
		return fmt.Sprintf("%s.%s", i.Task.Definition.Version, i.Task.Definition.Build)
	}

	return i.Task.Definition.Version
}

// Version ...
func (i Instance) Version() string {

	if i.Task.Definition.Version == "" {
		return "undetermined"
	}

	return i.Task.Definition.Version
}

// Build ...
func (i Instance) Build() string {

	if i.Task.Definition.Build == "" {
		return "undetermined"
	}

	return i.Task.Definition.Build
}
