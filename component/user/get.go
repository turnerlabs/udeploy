package user

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"gopkg.in/mgo.v2/bson"
)

// GetAll ...
func GetAll(ctx context.Context) ([]User, error) {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("users")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []User{}, err
	}
	defer cur.Close(ctx)

	users := []User{}
	for cur.Next(ctx) {
		user := &User{}

		if err := cur.Decode(user); err != nil {
			return []User{}, err
		}

		users = append(users, *user)
	}

	if err := cur.Err(); err != nil {
		return []User{}, err
	}

	return users, nil
}

// Get ...
func Get(ctx context.Context, user string) (User, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("users")

	match := bson.M{"email": user}

	usr := User{}
	if err := collection.FindOne(ctx, match).Decode(&usr); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return User{}, fmt.Errorf("%s is not a valid user", user)
		}

		return User{}, err
	}

	return usr, nil
}
