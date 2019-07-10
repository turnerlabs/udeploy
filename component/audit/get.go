package audit

import (
	"context"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetEntriesByAppInstance ...
func GetEntriesByAppInstance(ctx context.Context, app, instance string) ([]Entry, error) {
	collection := db.Client().Database(cfg.Get["DB_NAME"]).Collection("audit")

	match := bson.M{"app": app, "instance": instance}

	opts := options.FindOptions{
		Sort: bson.M{"time": -1},
	}

	cur, err := collection.Find(ctx, match, &opts)
	if err != nil {
		return []Entry{}, err
	}
	defer cur.Close(ctx)

	entries := []Entry{}

	for cur.Next(ctx) {
		entry := &Entry{}

		if err := cur.Decode(entry); err != nil {
			return []Entry{}, err
		}

		entries = append(entries, *entry)
	}

	if err := cur.Err(); err != nil {
		return []Entry{}, err
	}

	return entries, nil
}
