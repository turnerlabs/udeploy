package sync

import (
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/integration/aws/sqs"
	"github.com/turnerlabs/udeploy/component/supplement"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// AWSWatchEvents ...
func AWSWatchEvents(ctx mongo.SessionContext) error {
	return sqs.MonitorChanges(ctx, handleChange)
}

// AWSWatchS3 ...
func AWSWatchS3(ctx mongo.SessionContext) error {
	return sqs.MonitorS3(ctx, handleChange)
}

// AWSWatchAlarms ...
func AWSWatchAlarms(ctx mongo.SessionContext) error {
	return sqs.MonitorAlarms(ctx, handleChange)
}

func handleChange(ctx mongo.SessionContext, message sqs.MessageView) error {

	app, found := cache.Apps.GetByDefinitionID(message.ID)
	if !found {
		return nil
	}

	for name, inst := range app.Instances {
		if message.ID == inst.Task.Definition.ID {

			instances, err := supplement.Instances(ctx, app.Type, map[string]model.Instance{name: inst}, false)
			if err != nil {
				return err
			}

			cache.Apps.UpdateInstances(app.Name, instances)
		}
	}

	return nil
}
