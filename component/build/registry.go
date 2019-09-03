package build

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/integration/aws/s3"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/ecr"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"github.com/turnerlabs/udeploy/component/supplement"
)

type buildView struct {
	Revision   int64               `json:"revision"`
	Version    string              `json:"version"`
	Containers []app.ContainerView `json:"containers"`

	Registry bool `json:"registry"`
}

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
	appTypeS3            = "s3"
)

// GetBuilds ...
func GetBuilds(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	apps, err := app.Get(ctx, c.Param("app"))
	if err != nil {
		return err
	}

	instances := apps[0].GetInstances([]string{c.Param("registryInstance")})

	instances, err = supplement.Instances(ctx, apps[0].Type, instances, false)
	if err != nil {
		return err
	}

	sourceRegistry := instances[c.Param("registryInstance")]

	builds := map[string]app.Definition{}

	switch apps[0].Type {
	case appTypeService, appTypeScheduledTask:
		builds, err = task.ListDefinitions(sourceRegistry.Task)
		if err != nil {
			return err
		}

		if len(sourceRegistry.Repository) > 0 {
			ecrBuilds, err := ecr.ListDefinitions(sourceRegistry)
			if err != nil {
				return err
			}

			for _, b := range ecrBuilds {
				v := b.FormatVersion()

				if v == app.Undetermined {
					continue
				}

				if _, exists := builds[v]; !exists {
					builds[b.FormatVersion()] = b
				}
			}
		}

		// If the instance version no longer exists in the registry or
		// scanned task definitions, make it available for deployments.
		if _, exists := builds[sourceRegistry.FormatVersion()]; !exists {
			builds[sourceRegistry.FormatVersion()] = sourceRegistry.Task.Definition
		}

	case appTypeLambda:
		if len(sourceRegistry.S3RegistryBucket) > 0 {
			builds, err = s3.ListDefinitions(sourceRegistry)
			if err != nil {
				return err
			}
		} else {
			builds, err = lambda.ListDefinitions(sourceRegistry)
			if err != nil {
				return err
			}
		}
	case appTypeS3:
		builds, err = s3.ListDefinitions(sourceRegistry)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid app type %s", apps[0].Type)
	}

	viewBuilds := map[string]buildView{}

	for ver, details := range builds {
		revision := buildView{
			Revision: details.Revision,
			Version:  details.Version,
			Registry: details.Registry,
		}

		revision.Containers = append(revision.Containers, app.ContainerView{
			Image:       details.Description,
			Environment: details.Environment,
			Secrets:     details.Secrets,
		})

		viewBuilds[ver] = revision
	}

	return c.JSON(http.StatusOK, viewBuilds)
}
