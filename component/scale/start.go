package scale

import (
	"context"
	"fmt"
	"time"

	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/integration/aws/event"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/service"
	"github.com/turnerlabs/udeploy/model"
)

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
)

// Start ...
func Start(ctx context.Context, appType string, instance model.Instance, desiredCount int64, restart bool) error {
	switch appType {
	case appTypeService:
		return service.Scale(instance, desiredCount, restart)
	case appTypeScheduledTask:
		return event.Scale(ctx, instance, desiredCount, restart)
	case appTypeLambda:
		err := lambda.Scale(ctx, instance, desiredCount)

		// Currently udeploy does not watch the status of a running lambda
		// function. Pending is set to false after a few seconds to avoid
		// staying in the pending status for an indefinite time.
		go func() {
			time.Sleep(5 * time.Second)
			changeStatusToDone(instance)
		}()

		return err
	default:
		return fmt.Errorf("invalid app type %s", appType)
	}
}

func changeStatusToDone(instance model.Instance) {
	if app, found := cache.Apps.GetByDefinitionID(instance.Task.Definition.ID); found {
		for name, inst := range app.Instances {
			if inst.Task.Definition.ID == instance.Task.Definition.ID {
				inst.CurrentState.IsPending = false

				cache.Apps.UpdateInstances(app.Name, map[string]model.Instance{name: inst})
			}
		}
	}
}
