package handler

import (
	"log"
	"net/http"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/supplement"
	"github.com/turnerlabs/udeploy/component/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const appTypeLambda = "lambda"

// GetApp ..
func GetApp(c echo.Context) error {
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

// FilterApps ..
func FilterApps(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)
	usr := ctx.Value(session.ContextKey("user")).(user.User)

	filter := app.Filter{}
	if err := c.Bind(&filter); err != nil {
		return err
	}

	views := []app.AppView{}

	for name := range usr.Apps {

		if err := cache.EnsureApp(ctx, name); err != nil {
			c.Logger().Error(err)
			continue
		}

		application, _ := cache.Apps.Get(name)

		if application.Matches(filter) {
			views = append(views, application.ToView(usr))
		}
	}

	return c.JSON(http.StatusOK, views)
}

// GetApps ..
func GetApps(c echo.Context) error {
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

// SaveApp ...
func SaveApp(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)
	usr := ctx.Value(session.ContextKey("user")).(user.User)

	v := app.AppView{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	originalAppName := c.Param("app")
	originalApp, _ := cache.Apps.Get(originalAppName)

	newApp := v.ToBusiness()

	switch v.Type {
	case appTypeLambda:
		if arn, found := cfg.Get["SNS_ALARM_TOPIC_ARN"]; found {
			for _, i := range newApp.Instances {
				if err := lambda.UpsertAlarm(i.FunctionName, i.FunctionAlias, arn); err != nil {
					return err
				}
			}

			for name, i := range originalApp.Instances {
				if _, found := newApp.Instances[name]; !found {
					if err := lambda.DeleteAlarm(i.FunctionName); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}

	if err := app.Set(ctx, originalAppName, newApp); err != nil {
		return err
	}

	if originalAppName != v.Name {
		cache.Apps.Remove(originalAppName)

		usr.Apps[v.Name] = usr.Apps[originalAppName]
		delete(usr.Apps, originalAppName)
	}

	if !usr.Admin {
		app := usr.Apps[v.Name]
		for _, userInst := range v.Instances {
			if userInst.Edited {
				app.SetPermission(userInst.Name, "edit")
			}
		}
		usr.Apps[v.Name] = app
	}

	if err := user.Set(ctx, usr); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, v)
}

// DeleteApp ...
func DeleteApp(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	id, err := primitive.ObjectIDFromHex(c.Param("app"))
	if err != nil {
		return err
	}

	targetApp := cache.Apps.GetByID(id)

	switch targetApp.Type {
	case appTypeLambda:
		for _, i := range targetApp.Instances {
			if err := lambda.DeleteAlarm(i.FunctionName); err != nil {
				return err
			}
		}
	}

	if err := app.Delete(ctx, id); err != nil {
		return err
	}

	cache.Apps.RemoveByID(id)

	type response struct{}

	return c.JSON(http.StatusOK, response{})
}
