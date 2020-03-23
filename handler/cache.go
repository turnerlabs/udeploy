package handler

import (
	"net/http"

	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/project"
	"github.com/turnerlabs/udeploy/component/user"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/supplement"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetCachedApp ..
func GetCachedApp(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)
	usr := ctx.Value(session.ContextKey("user")).(user.User)

	apps, err := app.Get(ctx, c.Param("app"))
	if err != nil {
		return err
	}

	instances, err := supplement.Instances(ctx, apps[0].Type, apps[0].Instances, true)
	if err != nil {
		return err
	}

	apps[0].Instances = instances

	if err := cache.Apps.Update(apps[0]); err != nil {
		return err
	}

	project, err := project.Get(ctx, apps[0].ProjectID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, apps[0].ToView(usr, project))
}
