package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// RequireAdmin ...
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Get("ctx").(mongo.SessionContext)

		usr := ctx.Value(model.ContextKey("user")).(model.User)

		if usr.Admin {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("user %s unauthorized", usr.Email))
	}
}

// RequireDeploy ...
func RequireDeploy(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Get("ctx").(mongo.SessionContext)

		usr := ctx.Value(model.ContextKey("user")).(model.User)

		if usr.HasPermission(c.Param("app"), c.Param("instance"), "deploy") {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("user %s unauthorized for deploy on %s %s", usr.Email, c.Param("app"), c.Param("instance")))
	}
}

// RequireScale ...
func RequireScale(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Get("ctx").(mongo.SessionContext)

		usr := ctx.Value(model.ContextKey("user")).(model.User)

		if usr.HasPermission(c.Param("app"), c.Param("instance"), "scale") {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("user %s unauthorized for scale on %s %s", usr.Email, c.Param("app"), c.Param("instance")))
	}
}
