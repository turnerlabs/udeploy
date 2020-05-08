package s3

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Deploy ...
func Deploy(ctx mongo.SessionContext, actionID primitive.ObjectID, source app.Instance, target app.Instance, revision int64, opts task.DeployOptions) error {

	if len(opts.Secrets) > 0 {
		return errors.New("s3 does not support secrets")
	}

	session := session.New()

	config.Merge([]string{source.Role, source.RepositoryRole, target.Role}, session)

	workingDir := fmt.Sprintf("tmp/%s", uuid.New())
	if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
		return err
	}

	defer func() {
		if err := os.RemoveAll(workingDir); err != nil {
			log.Println(err)
		}
	}()

	version, err := getRevisionDetails(source, revision, session)
	if err != nil {
		return err
	}

	zipPath, err := download(source, revision, workingDir, session)
	if err != nil {
		return err
	}

	unzippedPath := fmt.Sprintf("%s/contents", workingDir)
	_, err = unzip(zipPath, unzippedPath)
	if err != nil {
		return err
	}

	metadata := map[string]*string{
		"version":  aws.String(version),
		"revision": aws.String(strconv.FormatInt(revision, 10)),
	}

	if !opts.Override {
		if opts.Environment, err = buildConfig(source, target, session); err != nil {
			return err
		}
	}

	if err := createConfigFile(unzippedPath, target, opts.Environment); err != nil {
		return err
	}

	if err = purge(target, session); err != nil {
		return err
	}

	if err = upload(target, unzippedPath, metadata, session); err != nil {
		return err
	}

	if len(target.CloudFrontID) > 0 {
		return invalidateCache(target, session)
	}

	return nil
}

func invalidateCache(target app.Instance, session *session.Session) error {
	now := time.Now().Format(time.RFC3339Nano)
	svc := cloudfront.New(session)
	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(target.CloudFrontID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			Paths: &cloudfront.Paths{
				Items:    aws.StringSlice([]string{"/*"}),
				Quantity: aws.Int64(1),
			},
			CallerReference: aws.String(now),
		},
	}

	_, err := svc.CreateInvalidation(input)
	return err
}

func createConfigFile(unzippedPath string, target app.Instance, env map[string]string) error {

	configPath := fmt.Sprintf("%s/%s", unzippedPath, target.S3ConfigKey)

	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(env); err != nil {
		return err
	}

	return nil
}

func getRevisionDetails(source app.Instance, revision int64, sess *session.Session) (string, error) {
	svc := s3.New(sess)

	o, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(source.S3RegistryBucket),
		Key:    aws.String(fmt.Sprintf("%s/%d.zip", source.S3RegistryPrefix, revision)),
	})
	if err != nil {
		return "", err
	}

	ver, _, err := extractMetadata(o.Metadata)
	if err != nil {
		return "", err
	}

	return ver, nil
}

func buildConfig(source, target app.Instance, sess *session.Session) (map[string]string, error) {

	if len(target.S3Prefix) > 0 {
		target.S3ConfigKey = fmt.Sprintf("%s/%s", target.S3Prefix, target.S3ConfigKey)
	}

	if len(source.S3Prefix) > 0 {
		source.S3ConfigKey = fmt.Sprintf("%s/%s", source.S3Prefix, source.S3ConfigKey)
	}

	buff := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buff).Encode(map[string]string{}); err != nil {
		return map[string]string{}, err
	}

	downloader := s3manager.NewDownloader(sess)

	dlBuff := aws.NewWriteAtBuffer(buff.Bytes())

	_, err := downloader.Download(dlBuff, &s3.GetObjectInput{
		Bucket: aws.String(target.S3Bucket),
		Key:    aws.String(target.S3ConfigKey),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
			default:
				return map[string]string{}, err
			}
		} else {
			return map[string]string{}, err
		}
	}

	newConfig := map[string]string{}
	if err := json.Unmarshal(dlBuff.Bytes(), &newConfig); err != nil {
		return newConfig, err
	}

	if len(target.Task.CloneEnvVars) > 0 {
		sourceBuff := aws.NewWriteAtBuffer([]byte{})

		_, err = downloader.Download(sourceBuff, &s3.GetObjectInput{
			Bucket: aws.String(source.S3Bucket),
			Key:    aws.String(source.S3ConfigKey),
		})
		if err != nil {
			return newConfig, err
		}

		sourceConfig := map[string]string{}

		if err := json.Unmarshal(sourceBuff.Bytes(), &sourceConfig); err != nil {
			return newConfig, err
		}

		for _, k := range target.Task.CloneEnvVars {
			if v, found := sourceConfig[k]; found {
				newConfig[k] = v
			}
		}
	}

	return newConfig, nil
}

func purge(target app.Instance, sess *session.Session) error {
	svc := s3.New(sess)

	listInput := &s3.ListObjectsInput{
		Bucket: aws.String(target.S3Bucket),
	}

	if len(target.S3Prefix) > 0 {
		listInput.SetPrefix(target.S3Prefix)
	}

	o, err := svc.ListObjects(listInput)
	if err != nil {
		return err
	}

	if len(o.Contents) == 0 {
		return nil
	}

	input := s3.DeleteObjectsInput{
		Bucket: aws.String(target.S3Bucket),
		Delete: &s3.Delete{},
	}

	objs := []*s3.ObjectIdentifier{}

	for _, obj := range o.Contents {
		if *obj.Key == target.S3FullConfigKey() {
			continue
		}

		objs = append(objs, &s3.ObjectIdentifier{
			Key: obj.Key,
		})
	}

	input.Delete.SetObjects(objs)

	_, err = svc.DeleteObjects(&input)
	if err != nil {
		return err
	}

	return nil
}

func upload(target app.Instance, workingDir string, metadata map[string]*string, sess *session.Session) error {
	uploader := s3manager.NewUploader(sess)

	objects := []s3manager.BatchUploadObject{}

	fileList := []string{}
	filepath.Walk(workingDir, func(path string, f os.FileInfo, err error) error {
		dir, err := isDirectory(path)
		if err != nil {
			return err
		}
		if dir {
			return nil
		}

		fileList = append(fileList, path)
		log.Printf("UPLOADING: %s\n", path)

		return nil
	})

	for _, filePath := range fileList {
		file, err := os.Open(filePath)
		if err != nil {
			log.Println("Failed to open file", filePath, err)
		}

		key := strings.Replace(filePath, workingDir, "", 1)
		if len(target.S3Prefix) > 0 {
			key = fmt.Sprintf("%s/%s", target.S3Prefix, key)
		}

		contentType := mime.TypeByExtension(filepath.Ext(filePath))

		objects = append(objects, s3manager.BatchUploadObject{
			Object: &s3manager.UploadInput{
				Key:         aws.String(key),
				Bucket:      aws.String(target.S3Bucket),
				Body:        file,
				Metadata:    metadata,
				ContentType: aws.String(contentType),
			},
		})

		defer file.Close()
	}

	iter := &s3manager.UploadObjectsIterator{Objects: objects}

	return uploader.UploadWithIterator(context.Background(), iter)
}

// download ...
func download(source app.Instance, revision int64, workingDir string, sess *session.Session) (string, error) {
	downloader := s3manager.NewDownloader(sess)

	zipFile := fmt.Sprintf("%s/deployment.zip", workingDir)

	file, err := os.Create(zipFile)
	if err != nil {
		return zipFile, err
	}
	defer file.Close()

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(source.S3RegistryBucket),
		Key:    aws.String(fmt.Sprintf("%s/%d.zip", source.S3RegistryPrefix, revision)),
	})
	if err != nil {
		return zipFile, err
	}

	return zipFile, nil
}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func isDirectory(path string) (bool, error) {
	fd, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}
