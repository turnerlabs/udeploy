package handler

import (
	"net/http"
	"strconv"

	"github.com/turnerlabs/udeploy/component/deploy"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeployRevision ...
func DeployRevision(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	opts := deploy.Options{}
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

	inst, err := deploy.Deploy(ctx, apps[0], c.Param("instance"), c.Param("registryInstance"), revision, opts)
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
