package cache

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/supplement"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// App ..
func App(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)
	usr := ctx.Value(model.ContextKey("user")).(model.User)

	apps, err := app.Get(ctx, c.Param("app"))
	if err != nil {
		return err
	}

	instances, err := supplement.Instances(ctx, apps[0].Type, apps[0].Instances, true)
	if err != nil {
		return err
	}

	apps[0].Instances = instances

	Apps.Update(apps[0])

	return c.JSON(http.StatusOK, apps[0].ToView(usr))
}
