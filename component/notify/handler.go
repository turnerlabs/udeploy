package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/model"
)

// Start ...
func Start(c echo.Context, updatesChan chan interface{}) error {
	ctx := c.Get("ctx").(context.Context)

	usr := ctx.Value(model.ContextKey("user")).(model.User)

	for {
		select {
		case <-c.Request().Context().Done():

			if c.Request().Context().Err() == context.Canceled {
				return nil
			}

			return c.Request().Context().Err()
		case msg, ok := <-updatesChan:
			if !ok {
				break
			}

			app, ok := msg.(model.Application)
			if !ok {
				log.Println("invalid message")
				continue
			}

			if _, isUserApp := usr.Apps[app.Name]; isUserApp {
				b, err := json.Marshal(app.ToView(usr))
				if err != nil {
					c.Logger().Debug(err)
				}

				fmt.Fprintf(c.Response().Writer, "data: %s\n\n", string(b))

				if f, ok := c.Response().Writer.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}
}
