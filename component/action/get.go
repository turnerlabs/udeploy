package action

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/bson"
)

// ErrNotFound ...
const ErrNotFound = "not found"

// Get ...
func Get(ctx context.Context, id primitive.ObjectID) (model.Action, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("actions")

	match := bson.M{"_id": id}

	action := model.Action{}
	if err := collection.FindOne(ctx, match).Decode(&action); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return model.Action{}, fmt.Errorf("%s is not a valid action", id)
		}

		return model.Action{}, err
	}

	return action, nil
}

// GetLatestBy ...
func GetLatestBy(ctx context.Context, definitionID string) (model.Action, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("actions")

	match := bson.M{"definitionId": definitionID}

	limit := int64(1.0)
	opts := options.FindOptions{
		Sort:  bson.M{"started": -1},
		Limit: &limit,
	}

	cur, err := collection.Find(ctx, match, &opts)
	if err != nil {
		return model.Action{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		a := model.Action{}

		if err := cur.Decode(&a); err != nil {
			return model.Action{}, err
		}

		return a, nil
	}

	return model.Action{}, errors.New(ErrNotFound)
}
