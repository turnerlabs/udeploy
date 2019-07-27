package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/broker"
	"github.com/turnerlabs/udeploy/component/notify"
)

// Change ...
func Change(c echo.Context, changeNotifier *broker.Broker) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")

	fmt.Fprint(c.Response().Writer, "event: open\n\n")
	if f, ok := c.Response().Writer.(http.Flusher); ok {
		f.Flush()
	}

	changes := changeNotifier.Subscribe()
	defer changeNotifier.Unsubscribe(changes)

	return notify.Start(c, changes)
}
