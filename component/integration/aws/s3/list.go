package s3

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/turnerlabs/udeploy/model"
)

const ext = ".zip"

// ListTaskDefinitions ...
func ListTaskDefinitions(instance model.Instance) (map[string]model.Definition, error) {

	results, err := List(instance.S3RegistryBucket, instance.S3RegistryPrefix)
	if err != nil {
		return nil, err
	}

	versions := map[string]model.Definition{}

	for _, o := range results {
		if !strings.Contains(*o.Key, ext) {
			continue
		}

		svc := s3.New(session.New())

		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(instance.S3RegistryBucket),
			Key:    aws.String(*o.Key),
		})
		if err != nil {
			log.Println(*o.Key)
			log.Println(err)
			continue
		}

		version, build, revision, err := extractMetadata(result.Metadata)
		if err != nil {
			log.Println(*o.Key)
			log.Println(err)
			continue
		}

		n, err := strconv.ParseInt(revision, 10, 64)
		if err != nil {
			log.Println(*o.Key)
			log.Println(err)
			continue
		}

		def := model.Definition{
			Version:  version,
			Build:    build,
			Revision: n,

			Environment: map[string]string{},
			Secrets:     map[string]string{},

			Description: fmt.Sprintf("%s.%s", version, build),
		}

		versions[def.FormatVersion()] = def
	}

	return versions, nil
}

// List ...
func List(bucket, registry string) ([]*s3.Object, error) {
	svc := s3.New(session.New())

	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	if len(registry) > 0 {
		input.SetPrefix(registry)
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		return []*s3.Object{}, err
	}

	return result.Contents, nil
}
