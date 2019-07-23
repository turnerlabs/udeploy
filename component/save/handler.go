package save

import (
	"log"
	"net/http"

	"github.com/turnerlabs/udeploy/component/user"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/notice"
	"go.mongodb.org/mongo-driver/mongo"
)

const appTypeLambda = "lambda"

// App ...
func App(c echo.Context) error {
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

// User ...
func User(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	usr := user.User{}
	if err := c.Bind(&usr); err != nil {
		return err
	}

	if err := user.Set(ctx, usr); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

// Notice ...
func Notice(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	n := notice.Notice{}
	if err := c.Bind(&n); err != nil {
		return err
	}

	if err := notice.Set(ctx, n); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, n)
}
