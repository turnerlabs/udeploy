package scale

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

// Instance ...
func Instance(c echo.Context, restart bool) error {
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

	if err := Start(ctx, app.Type, targetInstance, desiredCount, restart); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, success{})
}

type success struct{}
