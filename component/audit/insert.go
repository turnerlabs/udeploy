package audit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/turnerlabs/udeploy/component/user"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
)

// Entry ...
type Entry struct {
	App      string    `json:"app"`
	Instance string    `json:"instance"`
	Action   string    `json:"action"`
	User     string    `json:"user"`
	Time     time.Time `json:"time"`
}

// CreateEntry ...
func CreateEntry(ctx context.Context, app, instance, action string) error {
	e := Entry{
		Time:     time.Now().UTC(),
		App:      app,
		Instance: instance,
		Action:   action,
	}

	if usr, ok := ctx.Value(session.ContextKey("user")).(user.User); ok {
		e.User = usr.Email
	}

	if _, exists := cfg.Get["DB_URI"]; !exists {
		log.Println(e)
		return nil
	}

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("audit")

	if _, err := collection.InsertOne(ctx, e); err != nil {
		return fmt.Errorf("failed to audit %s for %s %s: %v", action, app, instance, err)
	}

	return nil
}
