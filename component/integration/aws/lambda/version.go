package lambda

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/version"
)

func extractVersion(instance app.Instance, config *lambda.FunctionConfiguration) (version.Version, error) {

	version, err := version.Extract(*config.Description, instance.Task.ImageTagEx)
	if err != nil {
		return version, err
	}

	if len(version.Build) == 0 {
		version.Build = (*config.RevisionId)[0:8]
	}

	return version, nil
}
