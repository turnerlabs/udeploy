package cache

import "github.com/labstack/echo/v4"

// NoCache ...
func NoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store")
		return next(c)
	}
}
