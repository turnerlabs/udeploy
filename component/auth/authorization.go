package auth

import (
	"fmt"
	"net/http"

	"github.com/turnerlabs/udeploy/component/user"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// RequireAdmin ...
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Get("ctx").(mongo.SessionContext)

		usr := ctx.Value(session.ContextKey("user")).(user.User)

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

		usr := ctx.Value(session.ContextKey("user")).(user.User)

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

		usr := ctx.Value(session.ContextKey("user")).(user.User)

		if usr.HasPermission(c.Param("app"), c.Param("instance"), "scale") {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("user %s unauthorized for scale on %s %s", usr.Email, c.Param("app"), c.Param("instance")))
	}
}

// RequireEdit ...
func RequireEdit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Get("ctx").(mongo.SessionContext)

		usr := ctx.Value(session.ContextKey("user")).(user.User)

		if usr.HasPermission(c.Param("app"), c.Param("instance"), "edit") {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("user %s unauthorized for edit on %s %s", usr.Email, c.Param("app"), c.Param("instance")))
	}
}
