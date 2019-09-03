package ecr

import (
	"strings"

	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/version"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ListDefinitions ...
func ListDefinitions(registry app.Instance) (map[string]app.Definition, error) {
	images, err := List(registry.Repo())
	if err != nil {
		return map[string]app.Definition{}, err
	}

	builds := map[string]app.Definition{}

	for _, i := range images {
		if i.ImageTag == nil {
			continue
		}

		ver, build := version.Extract(*i.ImageTag, registry.Task.ImageTagEx)

		builds[*i.ImageTag] = app.Definition{
			Version:  ver,
			Build:    build,
			Revision: registry.Task.Definition.Revision,

			Description: *i.ImageTag,

			Registry: true,
		}
	}

	return builds, nil
}

// List ...
func List(repo string) ([]*ecr.ImageIdentifier, error) {
	svc := ecr.New(session.New())

	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repo[strings.Index(repo, "/")+1 : len(repo)]),
		RegistryId:     aws.String(repo[0:strings.Index(repo, ".")]),
		Filter: &ecr.ListImagesFilter{
			TagStatus: aws.String(ecr.TagStatusTagged),
		},
	}

	result, err := svc.ListImages(input)
	if err != nil {
		return []*ecr.ImageIdentifier{}, err
	}

	return result.ImageIds, nil
}
