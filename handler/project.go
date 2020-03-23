package handler

import (
	"net/http"

	"github.com/turnerlabs/udeploy/component/project"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetProjects ...
func GetProjects(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	Projects, err := project.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Projects)
}

// SaveProject ...
func SaveProject(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	n := project.Project{}
	if err := c.Bind(&n); err != nil {
		return err
	}

	if err := project.Set(ctx, n); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, n)
}

// DeleteProject ...
func DeleteProject(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	if err := project.Delete(ctx, id); err != nil {
		return err
	}

	type response struct{}

	return c.JSON(http.StatusOK, response{})
}
