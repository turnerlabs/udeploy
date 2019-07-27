package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/turnerlabs/udeploy/component/scale"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

// Scale ...
func Scale(c echo.Context) error {
	return scaleInstance(c, false)
}

// Restart ...
func Restart(c echo.Context) error {
	return scaleInstance(c, true)
}

func scaleInstance(c echo.Context, restart bool) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	app, found := cache.Apps.Get(c.Param("app"))
	if !found {
		return fmt.Errorf("%s app not found", c.Param("app"))
	}

	desiredCount, err := strconv.ParseInt(c.Param("desiredCount"), 10, 64)
	if err != nil {
		return err
	}

	targetInstance, exists := app.Instances[c.Param("instance")]
	if !exists {
		return fmt.Errorf("%s instance not found", c.Param("instance"))
	}

	if err := scale.Start(ctx, app.Type, targetInstance, desiredCount, restart); err != nil {
		return err
	}

	type success struct{}

	return c.JSON(http.StatusOK, success{})
}
