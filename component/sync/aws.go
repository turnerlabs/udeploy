package sync

import (
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/integration/aws/sqs"
	"github.com/turnerlabs/udeploy/component/supplement"
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

	application, found := cache.Apps.GetByDefinitionID(message.ID)
	if !found {
		return nil
	}

	for name, inst := range application.Instances {
		if message.ID == inst.Task.Definition.ID {

			instances, err := supplement.Instances(ctx, application.Type, map[string]app.Instance{name: inst}, false)
			if err != nil {
				return err
			}

			cache.Apps.UpdateInstances(application.Name, instances)
		}
	}

	return nil
}
