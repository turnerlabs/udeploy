package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// Audit ...
func Audit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Get("ctx").(mongo.SessionContext)

		b, _ := ioutil.ReadAll(c.Request().Body)

		data := auditData{}
		json.Unmarshal(b, &data)

		// Restore the io.ReadCloser to its original state
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b))

		if err := next(c); err != nil {
			return err
		}

		action := c.Request().URL.Path

		if data != (auditData{}) {
			action = fmt.Sprintf("%s (%s)", action, data.Version)
		}

		if err := CreateEntry(ctx, c.Param("app"), c.Param("instance"), action); err != nil {
			return err
		}

		return nil
	}
}

type auditData struct {
	Version string `json:"version"`
}
