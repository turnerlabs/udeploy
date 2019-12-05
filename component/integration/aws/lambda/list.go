package lambda

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/config"
)

// ISO8601 time format
const ISO8601 = "2006-01-02T15:04:05-0700"

// ListDefinitions ...
func ListDefinitions(instance app.Instance) (map[string]app.Definition, error) {

	session := session.New()

	config.Merge([]string{instance.Role}, session)

	svc := lambda.New(session)
	ver, err := listVersions(svc, instance.FunctionName, "", []*lambda.FunctionConfiguration{})
	if err != nil {
		return nil, err
	}

	ver = limitRevisions(ver, instance.Task.Revisions)

	versions := map[string]app.Definition{}
	for _, funcVersion := range ver {

		revision, err := strconv.ParseInt(*funcVersion.Version, 10, 64)
		if err != nil {
			continue
		}

		version, build, err := extractVersion(instance, funcVersion)
		if err != nil {
			continue
		}

		env := map[string]string{}
		for k, v := range funcVersion.Environment.Variables {
			value := *v
			env[k] = value
		}

		def := app.Definition{
			Description: version,

			Version:  version,
			Build:    build,
			Revision: revision,

			Environment: env,
			Secrets:     map[string]string{},
		}

		versions[def.FormatVersion()] = def
	}

	return versions, nil
}

func listVersions(svc *lambda.Lambda, functionName, nextToken string, versions []*lambda.FunctionConfiguration) ([]*lambda.FunctionConfiguration, error) {

	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(functionName),
	}

	if len(nextToken) > 0 {
		input.SetMarker(nextToken)
	}

	output, err := svc.ListVersionsByFunction(input)
	if err != nil {
		return nil, err
	}

	versions = append(versions, output.Versions...)

	if output.NextMarker == nil || len(*output.NextMarker) == 0 {
		return versions, nil
	}

	return listVersions(svc, functionName, *output.NextMarker, versions)
}

func limitRevisions(objs []*lambda.FunctionConfiguration, limit int) []*lambda.FunctionConfiguration {

	if limit == 0 || len(objs) <= limit {
		return objs
	}

	sort.Slice(objs, func(i, j int) bool {

		iTime, err := time.Parse(ISO8601, *objs[i].LastModified)
		if err != nil {
			log.Println(err)
			return false
		}

		jTime, err := time.Parse(ISO8601, *objs[j].LastModified)
		if err != nil {
			log.Println(err)
			return false
		}

		return iTime.After(jTime)
	})

	return objs[0:limit]
}
