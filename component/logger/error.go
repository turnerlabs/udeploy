package logger

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// LogErrors ...
func LogErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if err := next(c); err != nil {
			log.Println(c.Request().URL.Path, err)
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return nil
	}
}
