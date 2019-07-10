package model

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/version"
)

// Definition ...
type Definition struct {
	ID string `json:"id"`

	Version  string `json:"version"`
	Build    string `json:"build"`
	Revision int64  `json:"revision"`

	Description string `json:"description"`

	Environment map[string]string `json:"environment"`
	Secrets     map[string]string `json:"secrets"`
}

// FormatVersion ...
func (d Definition) FormatVersion() string {

	if d.Version == "" {
		return "undetermined"
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
		ID: *td.Family,

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
