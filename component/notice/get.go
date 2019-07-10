package notice

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/model"
	"gopkg.in/mgo.v2/bson"
)

// Get ...
func Get(ctx context.Context, app string) ([]model.Notice, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("notices")

	match := bson.M{"$or": []bson.M{bson.M{"apps": bson.M{"name": app}}, bson.M{"apps": bson.M{"$size": 0}}}}

	cur, err := collection.Find(ctx, match)
	if err != nil {
		return []model.Notice{}, err
	}
	defer cur.Close(ctx)

	notifications := []model.Notice{}
	for cur.Next(ctx) {
		n := &model.Notice{}

		if err := cur.Decode(n); err != nil {
			return []model.Notice{}, err
		}

		notifications = append(notifications, *n)
	}

	if err := cur.Err(); err != nil {
		return []model.Notice{}, err
	}

	return notifications, nil
}

// GetAll ...
func GetAll(ctx context.Context) ([]model.Notice, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("notices")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []model.Notice{}, err
	}
	defer cur.Close(ctx)

	notifications := []model.Notice{}
	for cur.Next(ctx) {
		n := &model.Notice{}

		if err := cur.Decode(n); err != nil {
			return []model.Notice{}, err
		}

		notifications = append(notifications, *n)
	}

	if err := cur.Err(); err != nil {
		return []model.Notice{}, err
	}

	return notifications, nil
}
