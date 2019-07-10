package user

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Delete ...
func Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("users")

	match := bson.M{"_id": id}

	if _, err := collection.DeleteOne(ctx, match); err != nil {
		return err
	}

	return nil
}
