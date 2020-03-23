package project

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Delete ...
func Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("projects")

	match := bson.M{"_id": id}

	if _, err := collection.DeleteOne(ctx, match); err != nil {
		return err
	}

	return nil
}
