package app

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Get ...
func Get(ctx context.Context, name string) ([]model.Application, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("apps")

	match := bson.M{}

	if usr, ok := ctx.Value(model.ContextKey("user")).(model.User); ok {
		match = bson.M{"name": bson.M{"$in": usr.ListApps()}}
	}

	if len(name) > 0 {
		match = bson.M{"name": name}
	}

	cur, err := collection.Find(ctx, match)
	if err != nil {
		return []model.Application{}, err
	}
	defer cur.Close(ctx)

	apps := []model.Application{}
	for cur.Next(ctx) {
		app := &model.Application{}

		if err := cur.Decode(app); err != nil {
			return []model.Application{}, err
		}

		apps = append(apps, *app)
	}

	if err := cur.Err(); err != nil {
		return []model.Application{}, err
	}

	return apps, nil
}
