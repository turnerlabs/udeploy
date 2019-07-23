package lambda

import (
	"github.com/turnerlabs/udeploy/component/app"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/google/uuid"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
)

// Deploy ...
func Deploy(source app.Instance, target app.Instance, revision int64, opts task.DeployOptions) error {

	if opts.OverrideSecrets() {
		return errors.New("lambda functions do not support secrets")
	}

	svc := lambda.New(session.New())

	rev := strconv.FormatInt(revision, 10)

	sourceFuncArn := fmt.Sprintf("%s:%d", source.FunctionName, revision)

	if source.FunctionName != target.FunctionName || opts.Override() {

		sourceFunc, err := svc.GetFunction(&lambda.GetFunctionInput{
			FunctionName: aws.String(sourceFuncArn),
		})
		if err != nil {
			return err
		}

		if err := deployCode(*sourceFunc.Code.Location, target.FunctionName, svc); err != nil {
			return err
		}

		if err := deployConfig(target, sourceFunc.Configuration.Environment.Variables, opts, svc); err != nil {
			return err
		}

		vo, err := svc.PublishVersion(&lambda.PublishVersionInput{
			FunctionName: aws.String(target.FunctionName),
			Description:  sourceFunc.Configuration.Description,
		})
		if err != nil {
			return err
		}

		fmt.Println(*sourceFunc.Configuration.Description)

		lo, err := svc.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
			FunctionName: aws.String(target.FunctionName),
		})
		if err != nil {
			return err
		}

		for _, v := range lo.Versions {
			if *v.Description == *vo.Description && *v.Version != *vo.Version {
				_, err := svc.DeleteFunction(&lambda.DeleteFunctionInput{
					FunctionName: v.FunctionName,
					Qualifier:    v.Version,
				})
				if err != nil {
					log.Print(err)
				}
			}
		}

		rev = *vo.Version
	}

	fmt.Println(target.FunctionAlias)
	fmt.Println(target.FunctionName)
	fmt.Println(rev)

	_, err := svc.UpdateAlias(&lambda.UpdateAliasInput{
		Name:            aws.String(target.FunctionAlias),
		FunctionName:    aws.String(target.FunctionName),
		FunctionVersion: aws.String(rev),
	})
	if err != nil {
		return err
	}

	return nil
}

func deployConfig(target app.Instance, sourceEnvironment map[string]*string, opts task.DeployOptions, svc *lambda.Lambda) error {
	to, err := svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(target.FunctionName),
	})
	if err != nil {
		return err
	}

	input := mapConfiguration(*to.Configuration)

	for _, key := range target.Task.CloneEnvVars {
		if value, found := sourceEnvironment[key]; found {
			input.Environment.Variables[key] = value
		}
	}

	if opts.OverrideEnvironment() {
		input.Environment.SetVariables(aws.StringMap(opts.Environment))
	}

	_, err = svc.UpdateFunctionConfiguration(&input)
	if err != nil {
		return err
	}

	return nil
}

func deployCode(sourceCodeURL, targetFuncName string, svc *lambda.Lambda) error {
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
	if err != nil {
		return err
	}

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
		RevisionId:       config.RevisionId,
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
