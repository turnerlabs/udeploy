package scale

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/turnerlabs/udeploy/component/db"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/action"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/turnerlabs/udeploy/component/integration/aws/event"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/service"
)

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
)

// Start ...
func Start(ctx mongo.SessionContext, appType string, instance app.Instance, desiredCount int64, restart bool) error {
	switch appType {
	case appTypeService:
		_, err := action.Start(ctx, instance.Task.Definition.ID, "scale", time.Second*120)
		if err != nil {
			return err
		}

		return service.Scale(instance, desiredCount, restart)
	case appTypeScheduledTask:
		_, err := action.Start(ctx, instance.Task.Definition.ID, "scale", time.Second*120)
		if err != nil {
			return err
		}

		return event.Scale(ctx, instance, desiredCount, restart)
	case appTypeLambda:
		id, err := action.Start(ctx, instance.Task.Definition.ID, "scale", time.Second*30)
		if err != nil {
			return err
		}

		err = lambda.Scale(ctx, instance, desiredCount)

		// Currently udeploy does not watch the status of a running lambda
		// function. Ready is set to true after a few seconds to avoid
		// staying in the pending status for an indefinite time.
		go func(oid primitive.ObjectID) {
			time.Sleep(10 * time.Second)

			bgctx := context.Background()

			sess, err := db.Client().StartSession()
			if err != nil {
				log.Println(err)
				return
			}
			defer sess.EndSession(bgctx)

			if err := mongo.WithSession(bgctx, sess, func(sctx mongo.SessionContext) error {
				return action.Stop(sctx, oid, nil)
			}); err != nil {
				log.Printf("Failed to stop action: %s\n", oid)
				log.Println(err)
			}
		}(id)

		return err
	default:
		return fmt.Errorf("invalid app type %s", appType)
	}
}
