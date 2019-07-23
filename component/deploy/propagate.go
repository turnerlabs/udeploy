package deploy

import (
	"fmt"
	"log"

	"github.com/turnerlabs/udeploy/component/app"

	"github.com/turnerlabs/udeploy/component/audit"
	"go.mongodb.org/mongo-driver/mongo"
)

// Propagate ...
func Propagate(ctx mongo.SessionContext, messages chan interface{}) error {

	for msg := range messages {
		application, ok := msg.(app.Application)
		if !ok {
			log.Println("invalid message")
			continue
		}

		for name, inst := range application.Instances {

			if changed, changes := inst.Changed(); changed {
				if _, found := changes["VERSION"]; found {

					for targetName, target := range application.Instances {
						if target.AutoPropagate && target.Task.Registry == name && targetName != name {
							updatedInst, err := deploy(ctx, application, targetName, name, inst.Task.Definition.Revision, deployOptions{})
							if err != nil {
								log.Println(err)
								continue
							}

							action := fmt.Sprintf("automatic %s deployment trigger by %s deployment (%s)", targetName, name, updatedInst.FormatVersion())

							if err := audit.CreateEntry(ctx, application.Name, targetName, action); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}
