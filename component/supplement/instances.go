package supplement

import (
	"errors"
	"fmt"

	"github.com/turnerlabs/udeploy/component/action"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/integration/aws/s3"

	"github.com/turnerlabs/udeploy/component/integration/aws/event"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/service"
	"github.com/turnerlabs/udeploy/model"
)

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
	appTypeS3            = "s3"
)

// Instances ...
func Instances(ctx mongo.SessionContext, appType string, instances map[string]model.Instance, details bool) (insts map[string]model.Instance, err error) {
	switch appType {
	case appTypeService:
		insts, err = service.Populate(instances, details)
	case appTypeScheduledTask:
		insts, err = event.Populate(instances, details)
	case appTypeLambda:
		insts, err = lambda.Populate(instances, details)
	case appTypeS3:
		insts, err = s3.Populate(instances, details)
	default:
		return nil, fmt.Errorf("invalid app type %s", appType)
	}

	if err != nil {
		return insts, err
	}

	return checkCurrentActions(ctx, insts)

}

func checkCurrentActions(ctx mongo.SessionContext, instances map[string]model.Instance) (map[string]model.Instance, error) {

	for key, i := range instances {
		a, err := action.GetLatestBy(ctx, i.Task.Definition.ID)
		if err != nil {
			if err.Error() == action.ErrNotFound {
				continue
			}

			return instances, err
		}

		if a.Is(model.Pending) {
			i.CurrentState.IsPending = true
			i.CurrentState.IsRunning = false
		}

		if a.Is(model.Error) {
			if i.CurrentState.Error != nil {
				i.CurrentState.Error = fmt.Errorf("%s: %s", a.Info, i.CurrentState.Error)
			} else {
				i.CurrentState.Error = errors.New(a.Info)
			}

			i.CurrentState.IsRunning = false
		}

		instances[key] = i
	}

	return instances, nil
}
