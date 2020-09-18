package ecr

import (
	"strings"

	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/version"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ListDefinitions ...
func ListDefinitions(registry app.Instance) (map[string]app.Definition, error) {
	images, err := List(registry)
	if err != nil {
		return map[string]app.Definition{}, err
	}

	builds := map[string]app.Definition{}

	for _, i := range images {
		if i.ImageTag == nil {
			continue
		}

		ver, err := version.Extract(*i.ImageTag, registry.Task.ImageTagEx)
		if err != nil {
			continue
		}

		builds[*i.ImageTag] = app.Definition{
			Version:  ver,
			Revision: registry.Task.Definition.Revision,

			Description: *i.ImageTag,

			Registry: true,
		}
	}

	return builds, nil
}

// List ...
func List(i app.Instance) ([]*ecr.ImageIdentifier, error) {

	repo := i.Repo()

	session := session.New()

	if len(i.RepositoryRole) > 0 {
		session.Config.WithCredentials(stscreds.NewCredentials(session, i.RepositoryRole))
	}

	svc := ecr.New(session)

	imageIds, err := list(svc, repo, "", []*ecr.ImageIdentifier{})
	if err != nil {
		return []*ecr.ImageIdentifier{}, err
	}

	return imageIds, nil
}

func list(svc *ecr.ECR, repo, nextToken string, ids []*ecr.ImageIdentifier) ([]*ecr.ImageIdentifier, error) {

	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repo[strings.Index(repo, "/")+1 : len(repo)]),
		RegistryId:     aws.String(repo[0:strings.Index(repo, ".")]),
		Filter: &ecr.ListImagesFilter{
			TagStatus: aws.String(ecr.TagStatusTagged),
		},
	}

	if len(nextToken) > 0 {
		input.SetNextToken(nextToken)
	}

	output, err := svc.ListImages(input)
	if err != nil {
		return nil, err
	}

	ids = append(ids, output.ImageIds...)

	if output.NextToken == nil || len(*output.NextToken) == 0 {
		return ids, nil
	}

	return list(svc, repo, *output.NextToken, ids)
}
