package handler

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUsers ...
func GetUsers(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	users, err := user.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// SaveUser ...
func SaveUser(c echo.Context) error {
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

// DeleteUser ...
func DeleteUser(c echo.Context) error {

	ctx := c.Get("ctx").(mongo.SessionContext)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	if err := user.Delete(ctx, id); err != nil {
		return err
	}

	type response struct{}

	return c.JSON(http.StatusOK, response{})
}
