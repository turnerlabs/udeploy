package audit

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAuditEntries ..
func GetAuditEntries(c echo.Context) error {
	ctx := c.Get("ctx").(mongo.SessionContext)

	entries, err := GetEntriesByAppInstance(ctx, c.Param("app"), c.Param("instance"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entries)
}
