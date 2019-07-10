package notify

import (
	"github.com/turnerlabs/udeploy/component/notice"
	"fmt"
	"log"

	"github.com/turnerlabs/udeploy/component/cache"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/integration/aws/sns"
	"github.com/turnerlabs/udeploy/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// Watch ...
func Watch(ctx mongo.SessionContext, messages chan interface{}) error {

	for msg := range messages {
		app, ok := msg.(model.Application)
		if !ok {
			log.Println("invalid message")
			continue
		}

		for instanceName, inst := range app.Instances {

			if changed, changes := inst.Changed(); changed {

				for t, c := range changes {
					log.Printf("CHANGED: %s [%s] %s\n", inst.Task.Definition.ID, t, c)
				}

				notifications, err := notice.Get(ctx, app.Name)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, n := range notifications {

					if n.Enabled && n.Matches(instanceName, inst) {

						subject := fmt.Sprintf("%s: %s %s %s (%s)", n.Name, app.Name, instanceName, inst.FormatVersion(), inst.String())

						body := fmt.Sprintf("%s\n\n %s/apps/%s/instance/%s", subject, cfg.Get["URL"], app.Name, instanceName)
						if inst.CurrentState.Error != nil {
							body = fmt.Sprintf("%s \n\n %s", body, inst.CurrentState.Error)
						}

						if err := sns.Publish(subject, body, n.SNSArn); err != nil {
							log.Println(err)
						}
					}
				}

				cache.Apps.ResetChangeState(app.Name, instanceName)
			}
		}
	}

	return nil
}
