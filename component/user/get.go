package user

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/model"
	"gopkg.in/mgo.v2/bson"
)

// GetAll ...
func GetAll(ctx context.Context) ([]model.User, error) {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("users")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []model.User{}, err
	}
	defer cur.Close(ctx)

	users := []model.User{}
	for cur.Next(ctx) {
		user := &model.User{}

		if err := cur.Decode(user); err != nil {
			return []model.User{}, err
		}

		users = append(users, *user)
	}

	if err := cur.Err(); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

// Get ...
func Get(ctx context.Context, user string) (model.User, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("users")

	match := bson.M{"email": user}

	usr := model.User{}
	if err := collection.FindOne(ctx, match).Decode(&usr); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return model.User{}, fmt.Errorf("%s is not a valid user", user)
		}

		return model.User{}, err
	}

	return usr, nil
}
