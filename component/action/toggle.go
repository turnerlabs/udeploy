package action

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Start ...
func Start(ctx context.Context, definitionID, aType string) (primitive.ObjectID, error) {
	return Set(ctx, model.Action{
		DefinitionID: definitionID,

		Type:    aType,
		Status:  model.Pending,
		Started: time.Now().UTC(),
	})
}

// Stop ...
func Stop(ctx mongo.SessionContext, id primitive.ObjectID, actionErr error) error {

	a, err := Get(ctx, id)
	if err != nil {
		return err
	}

	a.Status = model.Complete

	if actionErr != nil {
		a.Status = model.Error
		a.Info = actionErr.Error()
	}

	a.Stopped = time.Now().UTC()

	_, err = Set(ctx, a)

	return err
}
