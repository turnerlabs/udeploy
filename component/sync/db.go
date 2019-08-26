package sync

import (
	"log"

	"github.com/turnerlabs/udeploy/component/action"
	"github.com/turnerlabs/udeploy/component/app"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/component/supplement"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/turnerlabs/udeploy/component/cfg"
)

// WatchDatabaseApps ...
func WatchDatabaseApps(ctx mongo.SessionContext) error {

	c := db.Client().Database(cfg.Get["DB_NAME"]).Collection("apps")

	var pipeline []interface{}
	updateLookup := options.UpdateLookup

	cur, err := c.Watch(ctx, pipeline, &options.ChangeStreamOptions{
		FullDocument: &updateLookup,
	})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		change := struct {
			DocumentKey struct {
				ID primitive.ObjectID `bson:"_id"`
			} `bson:"documentKey"`
			OperationType string          `bson:"operationType"`
			App           app.Application `bson:"fullDocument"`
		}{}

		if err := cur.Decode(&change); err != nil {
			log.Fatal(err)
		}

		switch change.OperationType {
		case "insert", "update":
			instances, err := supplement.Instances(ctx, change.App.Type, change.App.Instances, false)
			if err != nil {
				log.Println("cache update failed after database app change detected")
				log.Println(err)
				continue
			}

			change.App.Instances = instances

			cache.Apps.Update(change.App)
		case "delete":
			cache.Apps.RemoveByID(change.DocumentKey.ID)
		}
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

// WatchDatabaseActions ...
func WatchDatabaseActions(ctx mongo.SessionContext) error {

	c := db.Client().Database(cfg.Get["DB_NAME"]).Collection("actions")

	var pipeline []interface{}
	updateLookup := options.UpdateLookup

	cur, err := c.Watch(ctx, pipeline, &options.ChangeStreamOptions{
		FullDocument: &updateLookup,
	})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		change := struct {
			DocumentKey struct {
				ID primitive.ObjectID `bson:"_id"`
			} `bson:"documentKey"`
			OperationType string        `bson:"operationType"`
			Action        action.Action `bson:"fullDocument"`
		}{}

		if err := cur.Decode(&change); err != nil {
			log.Fatal(err)
		}

		switch change.OperationType {
		case "insert", "update":
			app, found := cache.Apps.GetByDefinitionID(change.Action.DefinitionID)
			if !found {
				continue
			}

			instances, err := supplement.Instances(ctx, app.Type, app.Instances, false)
			if err != nil {
				log.Println("cache update failed after database action change detected")
				log.Println(err)
				continue
			}

			app.Instances = instances

			cache.Apps.Update(app)
		}
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}
