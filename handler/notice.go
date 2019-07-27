package handler

import (
	"net/http"

	"github.com/turnerlabs/udeploy/component/notice"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetNotices ...
func GetNotices(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	notices, err := notice.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, notices)
}

// SaveNotice ...
func SaveNotice(c echo.Context) error {
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

// DeleteNotice ...
func DeleteNotice(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	if err := notice.Delete(ctx, id); err != nil {
		return err
	}

	type response struct{}

	return c.JSON(http.StatusOK, response{})
}
