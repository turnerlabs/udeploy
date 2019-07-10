package app

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Set ...
func Set(ctx context.Context, appName string, app model.Application) error {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("apps")

	match := bson.M{"name": appName}

	d, err := toDoc(app)
	if err != nil {
		return err
	}

	update := bson.M{"$set": d}

	upsert := true
	_, err = collection.UpdateOne(ctx, match, update, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return err
	}

	return nil
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
