package notify

import (
	"fmt"
	"log"

	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/notice"

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

					if !n.Enabled {
						continue
					}

					for t, change := range changes {

						if n.Matches(instanceName, t, change) {

							subject := fmt.Sprintf("%s: %s %s %s (%s)", n.Name, application.Name, instanceName, inst.FormatVersion(), displayStatus(t, change))

							body := fmt.Sprintf("%s\n\n %s/apps/%s/instance/%s", subject, cfg.Get["URL"], application.Name, instanceName)
							if inst.CurrentState.Error != nil {
								body = fmt.Sprintf("%s \n\n %s", body, inst.CurrentState.Error)
							}

							if err := sns.Publish(subject, body, n.SNSArn); err != nil {
								log.Println(err)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func displayStatus(t string, c app.Change) string {

	switch t {
	case app.ChangeTypeVersion:
		return app.Deployed
	case app.ChangeTypeStatus:
		return c.After
	case app.ChangeTypeError:
		return app.Error
	}

	return "unknown"
}
