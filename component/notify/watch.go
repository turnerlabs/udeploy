package notify

import (
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/notice"
	"fmt"
	"log"

	"github.com/turnerlabs/udeploy/component/cache"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/integration/aws/sns"
	"go.mongodb.org/mongo-driver/mongo"
)

// Watch ...
func Watch(ctx mongo.SessionContext, messages chan interface{}) error {

	for msg := range messages {
		application, ok := msg.(app.Application)
		if !ok {
			log.Println("invalid message")
			continue
		}

		for instanceName, inst := range application.Instances {

			if changed, changes := inst.Changed(); changed {

				for t, c := range changes {
					log.Printf("CHANGED: %s [%s] %s\n", inst.Task.Definition.ID, t, c)
				}

				notifications, err := notice.Get(ctx, application.Name)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, n := range notifications {

					if n.Enabled && n.Matches(instanceName, inst) {

						subject := fmt.Sprintf("%s: %s %s %s (%s)", n.Name, application.Name, instanceName, inst.FormatVersion(), inst.String())

						body := fmt.Sprintf("%s\n\n %s/apps/%s/instance/%s", subject, cfg.Get["URL"], application.Name, instanceName)
						if inst.CurrentState.Error != nil {
							body = fmt.Sprintf("%s \n\n %s", body, inst.CurrentState.Error)
						}

						if err := sns.Publish(subject, body, n.SNSArn); err != nil {
							log.Println(err)
						}
					}
				}

				cache.Apps.ResetChangeState(application.Name, instanceName)
			}
		}
	}

	return nil
}
