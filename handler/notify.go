package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/project"
	"github.com/turnerlabs/udeploy/component/user"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/labstack/echo/v4"
)

// NotifyClients ...
func NotifyClients(c echo.Context, updatesChan chan interface{}) error {
	ctx := c.Get("ctx").(context.Context)

	usr := ctx.Value(session.ContextKey("user")).(user.User)

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

			app, ok := msg.(app.Application)
			if !ok {
				log.Println("invalid message")
				continue
			}

			if _, isUserApp := usr.Apps[app.Name]; isUserApp {
				project, err := project.Get(ctx, app.ProjectID)
				if err != nil {
					return err
				}

				b, err := json.Marshal(app.ToView(usr, project))
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
