package supplement

import (
	"fmt"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/turnerlabs/udeploy/component/action"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/integration/aws/s3"

	"github.com/turnerlabs/udeploy/component/integration/aws/event"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/service"
)

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
	appTypeS3            = "s3"
)

// Instances ...
func Instances(ctx mongo.SessionContext, appType string, instances map[string]app.Instance, details bool) (insts map[string]app.Instance, err error) {
	switch appType {
	case appTypeService:
		insts, err = service.Populate(instances)
	case appTypeScheduledTask:
		insts, err = event.Populate(instances)
	case appTypeLambda:
		insts, err = lambda.Populate(instances)
	case appTypeS3:
		insts, err = s3.Populate(instances)
	default:
		return nil, fmt.Errorf("invalid app type %s", appType)
	}

	if err != nil {
		return insts, err
	}

	return checkCurrentActions(ctx, insts)
}

func checkCurrentActions(ctx mongo.SessionContext, instances map[string]app.Instance) (map[string]app.Instance, error) {

	for key, i := range instances {
		actn, err := action.GetCurrentBy(ctx, i.Task.Definition.ID)
		if err != nil {
			if err.Error() == action.ErrNotFound {
				continue
			}

			return instances, err
		}

		switch actn.Status {
		case action.Pending:
			if actn.TimedOut() {
				err := app.StatusError{
					Value: fmt.Sprintf("%s action failed to receive a timely response from AWS", actn.Type),
					Type:  app.ErrorTypeAction,
				}

				if err := action.Stop(ctx, actn.ID, err); err != nil {
					return instances, err
				}
			}

			i.CurrentState.SetPending()
		case action.Error:
			if i.CurrentState.Error != nil {
				i.CurrentState.SetError(app.StatusError{Type: app.ErrorTypeAction, Value: fmt.Sprintf("%s: %s", actn.Info, i.CurrentState.Error)})
			} else {
				i.CurrentState.SetError(app.StatusError{Type: app.ErrorTypeAction, Value: actn.Info})
			}
		}

		instances[key] = i
	}

	return instances, nil
}
