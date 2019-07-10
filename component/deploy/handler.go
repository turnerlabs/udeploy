package deploy

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"go.mongodb.org/mongo-driver/mongo"
)

type deployOptions struct {
	Env      map[string]string `json:"env"`
	Secrets  map[string]string `json:"secrets"`
	ImageTag string            `json:"imageTag"`
}

func (o deployOptions) ToBusiness(repository string) task.DeployOptions {
	m := task.DeployOptions{
		Environment: o.Env,
		Secrets:     o.Secrets,
	}

	if len(repository) > 0 && len(o.ImageTag) > 0 {
		m.Image = fmt.Sprintf("%s:%s", repository, o.ImageTag)
	}

	return m
}

// Revision ...
func Revision(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	opts := deployOptions{}
	if err := c.Bind(&opts); err != nil {
		return err
	}

	apps, err := app.Get(ctx, c.Param("app"))
	if err != nil {
		return err
	}

	revision, err := strconv.ParseInt(c.Param("revision"), 10, 64)
	if err != nil {
		return err
	}

	inst, err := deploy(ctx, apps[0], c.Param("instance"), c.Param("registryInstance"), revision, opts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, deploySuccess{
		Image:    inst.Task.Definition.Description,
		Version:  inst.FormatVersion(),
		Revision: inst.Task.Definition.Revision,
	})
}

type deploySuccess struct {
	Version  string `json:"version"`
	Revision int64  `json:"revision"`
	Image    string `json:"image"`
}
