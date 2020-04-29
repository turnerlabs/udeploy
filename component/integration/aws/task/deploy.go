package task

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
	"github.com/turnerlabs/udeploy/component/version"
)

// DeployOptions ...
type DeployOptions struct {
	Override bool

	Environment map[string]string
	Secrets     map[string]string

	Image string
}

// DeployImage ...
func (do DeployOptions) DeployImage() bool {
	return len(do.Image) > 0
}

// Deploy ...
func Deploy(source app.Instance, target app.Instance, sourceRevision int64, sourceVersion string, opts DeployOptions) (td *ecs.TaskDefinition, err error) {
	session := session.New()

	config.Merge([]string{source.Role, target.Role}, session)

	svc := ecs.New(session)

	sourceTaskArn := fmt.Sprintf("%s:%d", source.Task.Family, sourceRevision)
	sourceOutput, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(sourceTaskArn),
	})
	if err != nil {
		return nil, err
	}

	targetTaskArn := fmt.Sprintf("%s:%d", target.Task.Family, target.Task.Definition.Revision)
	targetOutput, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(targetTaskArn),
	})
	if err != nil {
		return nil, err
	}

	containerDefinitions := []*ecs.ContainerDefinition{}

	if len(targetOutput.TaskDefinition.ContainerDefinitions) == len(sourceOutput.TaskDefinition.ContainerDefinitions) {
		for i, targetContainer := range targetOutput.TaskDefinition.ContainerDefinitions {
			sourceContainer := sourceOutput.TaskDefinition.ContainerDefinitions[i]

			targetContainer = targetContainer.SetImage(*sourceContainer.Image)

			secrets := targetContainer.Secrets

			if opts.Override {
				secrets = []*ecs.Secret{}

				for n, v := range opts.Secrets {
					name := n
					value := v

					secrets = append(secrets, &ecs.Secret{
						Name:      &name,
						ValueFrom: &value,
					})
				}
			}

			environment := targetContainer.Environment

			if opts.DeployImage() {
				targetContainer = targetContainer.SetImage(opts.Image)

				verEnv, buildEnv := source.RepoVersion()
				verValue, err := version.Extract(opts.Image, source.Task.ImageTagEx)
				if err != nil {
					return nil, err
				}

				environment = setEnvironmentVar(environment, verEnv, verValue.Version)
				environment = setEnvironmentVar(environment, buildEnv, verValue.Build)
			} else {
				environment = cloneEnvironment(sourceContainer.Environment, environment, target.Task.CloneEnvVars)
			}

			if opts.Override {
				environment = newEnvironment(opts.Environment)
			}

			targetContainer = targetContainer.SetSecrets(secrets)
			targetContainer = targetContainer.SetEnvironment(environment)

			containerDefinitions = append(containerDefinitions, targetContainer)
		}
	} else {
		return nil, fmt.Errorf("number of container definitions not compatible")
	}

	newOutput, err := svc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		Family:                  targetOutput.TaskDefinition.Family,
		ContainerDefinitions:    containerDefinitions,
		Cpu:                     targetOutput.TaskDefinition.Cpu,
		Memory:                  targetOutput.TaskDefinition.Memory,
		ExecutionRoleArn:        targetOutput.TaskDefinition.ExecutionRoleArn,
		TaskRoleArn:             targetOutput.TaskDefinition.TaskRoleArn,
		NetworkMode:             targetOutput.TaskDefinition.NetworkMode,
		PlacementConstraints:    targetOutput.TaskDefinition.PlacementConstraints,
		RequiresCompatibilities: targetOutput.TaskDefinition.RequiresCompatibilities,
		Volumes:                 targetOutput.TaskDefinition.Volumes,
	})
	if err != nil {
		return nil, err
	}

	return newOutput.TaskDefinition, nil
}

func newEnvironment(newEnvironment map[string]string) []*ecs.KeyValuePair {
	env := []*ecs.KeyValuePair{}

	for n, v := range newEnvironment {
		name := n
		value := v

		env = append(env, &ecs.KeyValuePair{
			Name:  &name,
			Value: &value,
		})
	}

	return env
}

func cloneEnvironment(source, target []*ecs.KeyValuePair, varsToClone []string) []*ecs.KeyValuePair {
	environment := []*ecs.KeyValuePair{}

	for _, varToClone := range varsToClone {
		for _, source := range source {
			if *source.Name == varToClone {
				environment = append(environment, source)
			}
		}
	}

	for _, v := range target {
		shouldAppend := true
		for _, clonedVar := range varsToClone {
			if *v.Name == clonedVar {
				shouldAppend = false
			}
		}
		if shouldAppend {
			environment = append(environment, v)
		}
	}

	return environment
}

func setEnvironmentVar(environment []*ecs.KeyValuePair, name, value string) []*ecs.KeyValuePair {
	if len(name) == 0 {
		return environment
	}

	newEnvironment := []*ecs.KeyValuePair{
		&ecs.KeyValuePair{
			Name:  aws.String(name),
			Value: aws.String(value),
		},
	}

	for _, pair := range environment {
		if *pair.Name != name {
			newEnvironment = append(newEnvironment, pair)
		}
	}

	return newEnvironment
}
