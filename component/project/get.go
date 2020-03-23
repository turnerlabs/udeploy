package project

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Get ...
func Get(ctx context.Context, id primitive.ObjectID) (Project, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("projects")

	match := bson.M{"_id": id}

	proj := Project{}
	if err := collection.FindOne(ctx, match).Decode(&proj); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return Project{}, nil
		}

		return Project{}, err
	}

	return proj, nil
}

// GetAll ...
func GetAll(ctx context.Context) ([]Project, error) {

	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("projects")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []Project{}, err
	}
	defer cur.Close(ctx)

	projects := []Project{}
	for cur.Next(ctx) {
		n := &Project{}

		if err := cur.Decode(n); err != nil {
			return []Project{}, err
		}

		projects = append(projects, *n)
	}

	if err := cur.Err(); err != nil {
		return []Project{}, err
	}

	return projects, nil
}
