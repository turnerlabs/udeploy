package action

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Set ...
func Set(ctx context.Context, a Action) (primitive.ObjectID, error) {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("actions")

	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
	}

	match := bson.M{"_id": a.ID}

	d, err := toDoc(a)
	if err != nil {
		return a.ID, err
	}

	update := bson.M{"$set": d}

	upsert := true
	_, err = collection.UpdateOne(ctx, match, update, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return a.ID, err
	}

	return a.ID, nil
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
