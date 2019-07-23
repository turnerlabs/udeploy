package app

import (
	"github.com/turnerlabs/udeploy/component/user"
	"context"

	"github.com/turnerlabs/udeploy/component/session"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Get ...
func Get(ctx context.Context, name string) ([]Application, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("apps")

	match := bson.M{}

	if usr, ok := ctx.Value(session.ContextKey("user")).(user.User); ok {
		match = bson.M{"name": bson.M{"$in": usr.ListApps()}}
	}

	if len(name) > 0 {
		match = bson.M{"name": name}
	}

	cur, err := collection.Find(ctx, match)
	if err != nil {
		return []Application{}, err
	}
	defer cur.Close(ctx)

	apps := []Application{}
	for cur.Next(ctx) {
		app := &Application{}

		if err := cur.Decode(app); err != nil {
			return []Application{}, err
		}

		apps = append(apps, *app)
	}

	if err := cur.Err(); err != nil {
		return []Application{}, err
	}

	return apps, nil
}
