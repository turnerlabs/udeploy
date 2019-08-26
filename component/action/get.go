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
	"go.mongodb.org/mongo-driver/bson"
)

// ErrNotFound ...
const ErrNotFound = "not found"

// Get ...
func Get(ctx context.Context, id primitive.ObjectID) (Action, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("actions")

	match := bson.M{"_id": id}

	action := Action{}
	if err := collection.FindOne(ctx, match).Decode(&action); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return Action{}, fmt.Errorf("%s is not a valid action", id)
		}

		return Action{}, err
	}

	return action, nil
}

// GetCurrentBy ...
func GetCurrentBy(ctx context.Context, definitionID string) (Action, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("actions")

	match := bson.D{{"definitionId", definitionID}}

	limit := int64(1.0)
	opts := options.FindOptions{
		Sort:  bson.M{"started": -1},
		Limit: &limit,
	}

	cur, err := collection.Find(ctx, match, &opts)
	if err != nil {
		return Action{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		a := Action{}

		if err := cur.Decode(&a); err != nil {
			return Action{}, err
		}

		return a, nil
	}

	return Action{}, errors.New(ErrNotFound)
}
