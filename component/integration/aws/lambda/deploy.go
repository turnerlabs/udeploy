package lambda

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/turnerlabs/udeploy/component/version"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/google/uuid"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Deploy ...
func Deploy(source app.Instance, target app.Instance, revision int64, opts task.DeployOptions) error {

	if len(opts.Secrets) > 0 {
		return errors.New("lambda functions do not support secrets")
	}

	if len(source.S3RegistryBucket) > 0 {
		return deployFromS3(source, target, revision, opts)
	}

	return deployFromLambda(source, target, revision, opts)
}

func deployFromLambda(source, target app.Instance, revision int64, opts task.DeployOptions) error {

	session := session.New()

	sourceConfig := aws.NewConfig()
	if len(source.Role) > 0 {
		sourceConfig.WithCredentials(stscreds.NewCredentials(session, source.Role))
	}

	targetConfig := aws.NewConfig()
	if len(target.Role) > 0 {
		targetConfig.WithCredentials(stscreds.NewCredentials(session, target.Role))
	}

	sourceSVC := lambda.New(session, sourceConfig)
	targetSVC := lambda.New(session, targetConfig)

	sourceFuncArn := fmt.Sprintf("%s:%d", source.FunctionName, revision)

	sourceFunc, err := sourceSVC.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(sourceFuncArn),
	})
	if err != nil {
		return err
	}

	if err := deployCodeFromLambda(*sourceFunc.Code.Location, target.FunctionName, targetSVC); err != nil {
		return err
	}

	if err := deployConfig(source, target, opts, sourceSVC, targetSVC); err != nil {
		return err
	}

	vo, err := targetSVC.PublishVersion(&lambda.PublishVersionInput{
		FunctionName: aws.String(target.FunctionName),
		Description:  sourceFunc.Configuration.Description,
	})
	if err != nil {
		return err
	}

	_, err = targetSVC.UpdateAlias(&lambda.UpdateAliasInput{
		Name:            aws.String(target.FunctionAlias),
		FunctionName:    aws.String(target.FunctionName),
		FunctionVersion: vo.Version,
	})

	return err
}

func deployFromS3(source, target app.Instance, revision int64, opts task.DeployOptions) error {
	session := session.New()

	sourceConfig := aws.NewConfig()
	if len(source.Role) > 0 {
		sourceConfig.WithCredentials(stscreds.NewCredentials(session, source.Role))
	}

	targetConfig := aws.NewConfig()
	if len(target.Role) > 0 {
		targetConfig.WithCredentials(stscreds.NewCredentials(session, target.Role))
	}

	repoConfig := aws.NewConfig()
	if len(source.RepositoryRole) > 0 {
		repoConfig.WithCredentials(stscreds.NewCredentials(session, source.RepositoryRole))
	}

	key := fmt.Sprintf("%d.zip", revision)
	if len(source.S3RegistryPrefix) > 0 {
		key = fmt.Sprintf("%s/%d.zip", source.S3RegistryPrefix, revision)
	}

	sourceSVC := lambda.New(session, sourceConfig)
	targetSVC := lambda.New(session, targetConfig)

	_, err := targetSVC.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(target.FunctionName),
		S3Bucket:     aws.String(source.S3RegistryBucket),
		S3Key:        aws.String(key),
	})
	if err != nil {
		return err
	}

	if err := deployConfig(source, target, opts, sourceSVC, targetSVC); err != nil {
		return err
	}

	s3svc := s3.New(session, repoConfig)

	ver, err := getVersion(source.S3RegistryBucket, key, s3svc)
	if err != nil {
		return err
	}

	vo, err := targetSVC.PublishVersion(&lambda.PublishVersionInput{
		FunctionName: aws.String(target.FunctionName),
		Description:  aws.String(ver),
	})
	if err != nil {
		return err
	}

	_, err = targetSVC.UpdateAlias(&lambda.UpdateAliasInput{
		Name:            aws.String(target.FunctionAlias),
		FunctionName:    aws.String(target.FunctionName),
		FunctionVersion: vo.Version,
	})

	return err
}

func getVersion(bucket, key string, s3svc *s3.S3) (string, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	output, err := s3svc.GetObject(input)
	if err != nil {
		return "", err
	}

	ver, found := output.Metadata["Version"]
	if !found {
		return "", errors.New("failed to get version from s3 object metadata")
	}

	return *ver, nil
}

func isOldRevision(config *lambda.FunctionConfiguration, deployConfig *lambda.FunctionConfiguration, regex string) bool {
	if *config.Version == *deployConfig.Version {
		return false
	}

	oldVersion, err := version.Extract(*config.Description, regex)
	if err != nil {
		return false
	}

	newVersion, err := version.Extract(*deployConfig.Description, regex)
	if err != nil {
		return false
	}

	return len(oldVersion.Version) > 0 && newVersion.Version == oldVersion.Version && newVersion.Build == oldVersion.Build
}

func deleteOldRevisions(target app.Instance, deployVersion *lambda.FunctionConfiguration, svc *lambda.Lambda) error {
	lo, err := svc.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(target.FunctionName),
	})
	if err != nil {
		return err
	}

	for _, v := range lo.Versions {
		if isOldRevision(v, deployVersion, target.Task.ImageTagEx) {
			_, err := svc.DeleteFunction(&lambda.DeleteFunctionInput{
				FunctionName: v.FunctionName,
				Qualifier:    v.Version,
			})
			if err != nil {
				log.Print(err)
			}
		}
	}

	return nil
}

func deployConfig(source, target app.Instance, opts task.DeployOptions, sourceSVC *lambda.Lambda, targetSVC *lambda.Lambda) error {

	sourceFunc, err := sourceSVC.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(fmt.Sprintf("%s:%s", source.FunctionName, source.FunctionAlias)),
	})
	if err != nil {
		return err
	}

	currentFunc, err := targetSVC.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(fmt.Sprintf("%s:%s", target.FunctionName, target.FunctionAlias)),
	})
	if err != nil {
		return err
	}

	input := mapConfiguration(*currentFunc.Configuration)

	for _, key := range target.Task.CloneEnvVars {
		if sourceFunc.Configuration.Environment != nil {
			if value, found := sourceFunc.Configuration.Environment.Variables[key]; found {
				input.Environment.Variables[key] = value
			}
		}
	}

	if opts.Override {
		input.Environment = &lambda.Environment{}
		input.Environment.SetVariables(aws.StringMap(opts.Environment))
	}

	_, err = targetSVC.UpdateFunctionConfiguration(&input)

	return err
}

func deployCodeFromLambda(sourceCodeURL, targetFuncName string, svc *lambda.Lambda) error {
	codeZipfileName := fmt.Sprintf("%s.zip", uuid.New())

	if err := downloadFile(codeZipfileName, sourceCodeURL); err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(codeZipfileName); err != nil {
			log.Println(err)
		}
	}()

	b, err := ioutil.ReadFile(codeZipfileName)
	if err != nil {
		return err
	}

	_, err = svc.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(targetFuncName),
		ZipFile:      b,
	})

	return nil
}

func mapConfiguration(config lambda.FunctionConfiguration) lambda.UpdateFunctionConfigurationInput {

	newConfig := lambda.UpdateFunctionConfigurationInput{
		DeadLetterConfig: config.DeadLetterConfig,
		Description:      config.Description,
		FunctionName:     config.FunctionName,
		Handler:          config.Handler,
		KMSKeyArn:        config.KMSKeyArn,
		MemorySize:       config.MemorySize,
		Role:             config.Role,
		Runtime:          config.Runtime,
		Timeout:          config.Timeout,
	}

	if config.TracingConfig != nil {
		newConfig.TracingConfig = &lambda.TracingConfig{
			Mode: config.TracingConfig.Mode,
		}
	}

	if config.Environment != nil {
		newConfig.Environment = &lambda.Environment{
			Variables: config.Environment.Variables,
		}
	}

	if config.VpcConfig != nil {
		newConfig.VpcConfig = &lambda.VpcConfig{
			SecurityGroupIds: config.VpcConfig.SecurityGroupIds,
			SubnetIds:        config.VpcConfig.SubnetIds,
		}
	}

	for _, l := range config.Layers {
		newConfig.Layers = append(newConfig.Layers, l.Arn)
	}

	return newConfig
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

func downloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
