package delete

import (
	"github.com/turnerlabs/udeploy/component/notice"
	"net/http"

	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"

	"github.com/turnerlabs/udeploy/component/cache"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const appTypeLambda = "lambda"

// App ...
func App(c echo.Context) error {
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

	return c.JSON(http.StatusOK, response{})
}

// User ...
func User(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	if err := user.Delete(ctx, id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response{})
}

// Notice ...
func Notice(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	if err := notice.Delete(ctx, id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response{})
}

type response struct{}
