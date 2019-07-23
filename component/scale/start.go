package scale

import (
	"context"
	"fmt"
	"time"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/turnerlabs/udeploy/component/cache"
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
func Start(ctx context.Context, appType string, instance app.Instance, desiredCount int64, restart bool) error {
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

func changeStatusToDone(instance app.Instance) {
	if application, found := cache.Apps.GetByDefinitionID(instance.Task.Definition.ID); found {
		for name, inst := range application.Instances {
			if inst.Task.Definition.ID == instance.Task.Definition.ID {
				inst.CurrentState.IsPending = false

				cache.Apps.UpdateInstances(application.Name, map[string]app.Instance{name: inst})
			}
		}
	}
}
