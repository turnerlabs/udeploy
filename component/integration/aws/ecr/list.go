package ecr

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

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
