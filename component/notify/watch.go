package notify

import (
	"fmt"
	"log"
	"strings"

	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/notice"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/integration/aws/sns"
	"github.com/turnerlabs/udeploy/component/integration/slack"
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

				logChanges(inst.Task.Definition.ID, changes)

				if configException(changes) {
					continue
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

							status := displayStatus(t, change)
							portalURL := fmt.Sprintf("%s/apps/%s/instance/%s", cfg.Get["URL"], application.Name, instanceName)

							switch {
							case n.IsSNS():
								subject := fmt.Sprintf("%s: %s %s %s (%s)", n.Name, application.Name, instanceName, inst.Task.Definition.Version.Full(), status)

								body := fmt.Sprintf("%s\n\n %s", subject, portalURL)
								if inst.CurrentState.Error != nil {
									body = fmt.Sprintf("%s \n\n %s", body, inst.CurrentState.Error)
								}

								if err := sns.Publish(subject, body, n.SNSArn); err != nil {
									log.Println(n.SNSArn)
									log.Println(err)
								}
							case n.IsSlack():
								subject := fmt.Sprintf(":%s: %s (%s)", Emoji(status), n.Name, strings.Title(status))

								body := fmt.Sprintf("App: %s\nInstance: %s\nVersion: %s", application.Name, instanceName, inst.Task.Definition.Version.Full())
								if inst.CurrentState.Error != nil {
									body = fmt.Sprintf("%s\nError: %s", body, inst.CurrentState.Error)
								}

								msg := slack.Template(subject, body, Color(status), "View Details", portalURL)

								if err := slack.Send(msg, n.SlackWebHook); err != nil {
									log.Println(n.SlackWebHook)
									log.Println(err)
								}
							}

						}
					}
				}
			}
		}
	}

	return nil
}

func configException(changes map[string]app.Change) bool {
	for t := range changes {
		if t == app.ChangeTypeException {
			return true
		}
	}

	return false
}

func logChanges(id string, changes map[string]app.Change) {
	for t, c := range changes {
		log.Printf("APP_CHANGED: %s [%s] %s\n", id, t, c)
	}
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
