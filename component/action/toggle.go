package action

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Start ...
func Start(ctx context.Context, definitionID, aType string, maxRuntime time.Duration) (primitive.ObjectID, error) {
	return Set(ctx, Action{
		DefinitionID: definitionID,

		Type:    aType,
		Status:  Pending,
		Started: time.Now().UTC(),

		Expiration: time.Now().Add(maxRuntime).UTC(),
	})
}

// Stop ...
func Stop(ctx mongo.SessionContext, id primitive.ObjectID, actionErr error) error {

	a, err := Get(ctx, id)
	if err != nil {
		return err
	}

	a.Status = Stopped

	if actionErr != nil {
		a.Status = Error
		a.Info = actionErr.Error()
	}

	a.Stopped = time.Now().UTC()

	_, err = Set(ctx, a)

	return err
}
