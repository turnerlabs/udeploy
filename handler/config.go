package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/integration/aws/secretsmanager"
)

// GetConfigValue ..
func GetConfigValue(c echo.Context) error {

	app, found := cache.Apps.Get(c.Param("app"))
	if !found {
		return fmt.Errorf("%s app not found", c.Param("app"))
	}

	inst, exists := app.Instances[c.Param("instance")]
	if !exists {
		return fmt.Errorf("%s instance not found", c.Param("instance"))
	}

	v, err := secretsmanager.GetForInstance(c.QueryParam("linkId"), inst)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ConfigLinkValue{
		LinkID: c.QueryParam("linkId"),
		Value:  v,
	})
}

// SaveConfigValue ..
func SaveConfigValue(c echo.Context) error {

	app, found := cache.Apps.Get(c.Param("app"))
	if !found {
		return fmt.Errorf("%s app not found", c.Param("app"))
	}

	inst, exists := app.Instances[c.Param("instance")]
	if !exists {
		return fmt.Errorf("%s instance not found", c.Param("instance"))
	}

	cl := ConfigLinkValue{}

	if err := c.Bind(&cl); err != nil {
		return err
	}

	if err := secretsmanager.UpdateForInstance(cl.LinkID, cl.Value, inst); err != nil {
		return err
	}

	type success struct{}

	return c.JSON(http.StatusOK, success{})
}

type ConfigLinkValue struct {
	LinkID string `json:"linkId"`
	Value  string `json:"value"`
}
