package s3

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/turnerlabs/udeploy/component/integration/aws/config"
	"github.com/turnerlabs/udeploy/component/version"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Populate ...
func Populate(instances map[string]app.Instance) (map[string]app.Instance, error) {

	for key, instance := range instances {
		sess := session.New()

		config.Merge([]string{instance.Role}, sess)

		svc := s3.New(sess)
		downloader := s3manager.NewDownloader(sess)

		i, state, err := populateInst(instance, svc, downloader)
		if err != nil {
			state.Error = err
		}

		i.SetState(state)

		instances[key] = i
	}

	return instances, nil
}

func populateInst(i app.Instance, svc *s3.S3, downloader *s3manager.Downloader) (app.Instance, app.State, error) {
	state := app.NewState()

	id := fmt.Sprintf("%s-%s", i.S3Bucket, i.S3ConfigKey)
	configKey := i.S3ConfigKey
	root := i.S3Bucket

	if len(i.S3Prefix) > 0 {
		id = fmt.Sprintf("%s-%s/%s", i.S3Bucket, i.S3Prefix, i.S3ConfigKey)
		configKey = fmt.Sprintf("%s/%s", i.S3Prefix, i.S3ConfigKey)
		root = fmt.Sprintf("%s/%s", i.S3Bucket, i.S3Prefix)
	}

	i.Task.Definition = app.NewDefinition(id)
	i.Task.DesiredCount = 1

	oo, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(i.S3Bucket),
		Key:    aws.String(configKey),
	})
	if err != nil {
		return i, state, err
	}

	ver, revision, err := extractMetadata(oo.Metadata)
	if err != nil {
		return i, state, err
	}

	v, err := version.Extract(ver, i.Task.ImageTagEx)
	if err != nil {
		state.SetError(err)
	}

	i.Task.Definition.Description = ver
	i.Task.Definition.Version = v

	i.Task.Definition.Revision, err = strconv.ParseInt(revision, 10, 64)
	if err != nil {
		return i, state, err
	}

	state.Version = i.Task.Definition.Version.Full()

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

func extractMetadata(metadata map[string]*string) (string, string, error) {
	version, found := metadata["Version"]
	if !found {
		return "", "", errors.New("version not found")
	}

	revision, found := metadata["Revision"]
	if !found {
		return "", "", errors.New("revision not found")
	}

	return *version, *revision, nil
}
