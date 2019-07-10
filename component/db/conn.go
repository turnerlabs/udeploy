package db // import "github.com/turnerlabs/udeploy/component/db"

import (
	"context"
	"errors"
	"log"

	"github.com/turnerlabs/udeploy/component/cfg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Client ...
func Client() *mongo.Client {
	if client == nil {
		log.Fatal(errors.New("database connection unintialized"))
	}

	return client
}

// Connect ...
func Connect(ctx context.Context, uri string) (err error) {

	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Connect(ctx); err != nil {
		return err
	}

	if _, exists := cfg.Get["DB_URI"]; exists {
		if err := client.Ping(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}

// Disconnect ...
func Disconnect(ctx context.Context) (err error) {
	return client.Disconnect(ctx)
}
