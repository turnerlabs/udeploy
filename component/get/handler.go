package get

import (
	"net/http"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/turnerlabs/udeploy/component/notice"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/supplement"
	"github.com/turnerlabs/udeploy/component/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// App ..
func App(c echo.Context) error {
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

	return c.JSON(http.StatusOK, apps[0].ToView(usr))
}

// Apps ..
func Apps(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)
	usr := ctx.Value(session.ContextKey("user")).(user.User)

	appNames := []string{}

	if usr.Admin && c.QueryParam("all") == "true" {
		for name := range cache.Apps.GetAll() {
			appNames = append(appNames, name)
		}
	} else {
		for name := range usr.Apps {
			appNames = append(appNames, name)
		}
	}

	views := []app.AppView{}

	for _, name := range appNames {

		if err := cache.EnsureApp(ctx, name); err != nil {
			c.Logger().Error(err)
			continue
		}

		application, _ := cache.Apps.Get(name)

		views = append(views, application.ToView(usr))
	}

	return c.JSON(http.StatusOK, views)
}

// Users ...
func Users(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	users, err := user.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// Notices ...
func Notices(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	notices, err := notice.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, notices)
}
