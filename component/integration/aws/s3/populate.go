package s3

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Populate ...
func Populate(instances map[string]app.Instance, details bool) (map[string]app.Instance, error) {

	sess := session.New()

	svc := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	for key, instance := range instances {

		i, state, err := populateInst(instance, svc, downloader)

		if err != nil {
			state.Error = err
		}

		state.Version = i.FormatVersion()

		i.SetState(state)

		instances[key] = i
	}

	return instances, nil
}

func populateInst(i app.Instance, svc *s3.S3, downloader *s3manager.Downloader) (app.Instance, app.State, error) {
	state := app.State{
		DesiredCount: 1,
	}

	id := fmt.Sprintf("%s-%s", i.S3Bucket, i.S3ConfigKey)
	configKey := i.S3ConfigKey
	root := i.S3Bucket

	if len(i.S3Prefix) > 0 {
		id = fmt.Sprintf("%s-%s/%s", i.S3Bucket, i.S3Prefix, i.S3ConfigKey)
		configKey = fmt.Sprintf("%s/%s", i.S3Prefix, i.S3ConfigKey)
		root = fmt.Sprintf("%s/%s", i.S3Bucket, i.S3Prefix)
	}

	i.Task.Definition = app.NewDefinition(id)

	oo, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(i.S3Bucket),
		Key:    aws.String(configKey),
	})
	if err != nil {
		return i, state, err
	}

	version, build, revision, err := extractMetadata(oo.Metadata)
	if err != nil {
		return i, state, err
	}

	i.Task.Definition.Description = fmt.Sprintf("%s.%s", version, build)
	i.Task.Definition.Version = version
	i.Task.Definition.Build = build

	i.Task.Definition.Revision, err = strconv.ParseInt(revision, 10, 64)
	if err != nil {
		return i, state, err
	}

	buff := &aws.WriteAtBuffer{}

	_, err = downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(i.S3Bucket),
		Key:    aws.String(configKey),
	})
	if err != nil {
		return i, state, err
	}

	if err := json.Unmarshal(buff.Bytes(), &i.Task.Definition.Environment); err != nil {
		return i, state, err
	}

	i.Links = append(i.Links, app.Link{
		Generated:   true,
		Description: "Root Path",
		Name:        "S3",
		URL: fmt.Sprintf("https://s3.console.aws.amazon.com/s3/buckets/%s/",
			root),
	})

	return i, state, nil
}

func extractMetadata(metadata map[string]*string) (string, string, string, error) {
	version, found := metadata["Version"]
	if !found {
		return "", "", "", errors.New("version not found")
	}

	build, found := metadata["Build"]
	if !found {
		return "", "", "", errors.New("build number not found")
	}

	revision, found := metadata["Revision"]
	if !found {
		return "", "", "", errors.New("revision not found")
	}

	return *version, *build, *revision, nil
}
