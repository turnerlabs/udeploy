package user

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Set ...
func Set(ctx context.Context, user model.User) error {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("users")

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	match := bson.M{"_id": user.ID}

	d, err := toDoc(user)
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
