package build

import (
	"fmt"
	"net/http"

	"github.com/turnerlabs/udeploy/component/version"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/integration/aws/s3"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/ecr"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"github.com/turnerlabs/udeploy/component/supplement"
	"github.com/turnerlabs/udeploy/model"
)

const (
	buildTypeImage    = "BUILD_TYPE_IMAGE"
	buildTypeRevision = "BUILD_TYPE_REVISION"
)

type buildView struct {
	Type       string                `json:"type"`
	Revision   int64                 `json:"revision"`
	Version    string                `json:"version"`
	Containers []model.ContainerView `json:"containers"`
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

	builds := map[string]model.Definition{}

	switch apps[0].Type {
	case appTypeService, appTypeScheduledTask:
		builds, err = task.ListDefinitions(sourceRegistry.Task)
		if err != nil {
			return err
		}
	case appTypeLambda:
		builds, err = lambda.ListDefinitions(sourceRegistry)
		if err != nil {
			return err
		}
	case appTypeS3:
		builds, err = s3.ListTaskDefinitions(sourceRegistry)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid app type %s", apps[0].Type)
	}

	viewBuilds := map[string]buildView{}

	for ver, details := range builds {
		revision := buildView{
			Type:     buildTypeRevision,
			Revision: details.Revision,
			Version:  details.Version,
		}

		revision.Containers = append(revision.Containers, model.ContainerView{
			Image:       details.Description,
			Environment: details.Environment,
			Secrets:     details.Secrets,
		})

		viewBuilds[ver] = revision
	}

	if len(sourceRegistry.Repository) > 0 {

		images, err := ecr.List(sourceRegistry.Repo())
		if err != nil {
			return err
		}

		for _, i := range images {
			if i.ImageTag == nil {
				continue
			}

			if _, exists := viewBuilds[*i.ImageTag]; exists {
				continue
			}

			ver, _ := version.Extract(*i.ImageTag, sourceRegistry.Task.ImageTagEx)
			
			viewBuilds[*i.ImageTag] = buildView{
				Type:    buildTypeImage,
				Version: ver,
				Containers: []model.ContainerView{
					model.ContainerView{
						Image: *i.ImageTag,
					},
				},
			}
		}
	}

	return c.JSON(http.StatusOK, viewBuilds)
}
