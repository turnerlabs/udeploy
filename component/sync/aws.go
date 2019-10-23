package sync

import (
	"log"
	"time"

	"github.com/turnerlabs/udeploy/component/action"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/integration/aws/sqs"
	"github.com/turnerlabs/udeploy/component/supplement"
	"go.mongodb.org/mongo-driver/mongo"
)

// AWSPollEventlessChanges monitors and updates the instance state
// for specific changes that do not trigger AWS events. This polling
// technique is only used when changes cannot be detected by AWS
// events. To avoid AWS rate limits, this technique should be used
// SPARINGLY.
//
// Currently AWS does not fire events when errors expire from ECS Task
// history. Since monitored errors in ECS Task history cause an a
// pplication to display an error state, the history must be monitored
// to determine when a service returns to a healthly state.
func AWSPollEventlessChanges(ctx mongo.SessionContext) error {

	ticker := time.NewTicker(app.FailedTaskExpiration * time.Minute)

	for {
		select {
		case <-ticker.C:
			for _, a := range cache.Apps.GetAll() {
				if a.Type != app.AppTypeService {
					continue
				}

				targeted := a.GetErrorInstances()
				if len(targeted) == 0 {
					continue
				}

				log.Printf("Updating App: %s\n", a.Name)

				supplemented, err := supplement.Instances(ctx, a.Type, targeted, false)
				if err != nil {
					log.Printf("failed to update %s state (%s)\n", a.Name, err)
					continue
				}

				cache.Apps.UpdateInstances(a.Name, supplemented)

				time.Sleep(time.Second)
			}
		}
	}
}

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

			if act, err := action.GetCurrentBy(ctx, inst.Task.Definition.ID); err == nil {

				if act.Is(action.Pending) || act.Is(action.Error) {
					if err := action.Stop(ctx, act.ID, nil); err != nil {
						return err
					}
				}

			} else if err.Error() != action.ErrNotFound {
				return err
			}

			instances, err := supplement.Instances(ctx, application.Type, map[string]app.Instance{name: inst}, false)
			if err != nil {
				return err
			}

			cache.Apps.UpdateInstances(application.Name, instances)
		}
	}

	return nil
}
