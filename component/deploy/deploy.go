package deploy

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/turnerlabs/udeploy/component/action"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/component/integration/aws/s3"
	"github.com/turnerlabs/udeploy/component/supplement"

	"github.com/turnerlabs/udeploy/component/integration/aws/event"
	"github.com/turnerlabs/udeploy/component/integration/aws/lambda"
	"github.com/turnerlabs/udeploy/component/integration/aws/service"
	"github.com/turnerlabs/udeploy/component/integration/aws/task"
	"github.com/turnerlabs/udeploy/model"
)

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
	appTypeS3            = "s3"
)

func deploy(ctx mongo.SessionContext, app model.Application, target, source string, revision int64, opts deployOptions) (model.Instance, error) {

	instances := app.GetInstances([]string{target, source})

	instances, err := supplement.Instances(ctx, app.Type, instances, false)
	if err != nil {
		return model.Instance{}, err
	}

	targetInstance, targetExists := instances[target]
	if !targetExists {
		return model.Instance{}, err
	}

	registryInstance, sourceExists := instances[source]
	if !sourceExists {
		return model.Instance{}, err
	}

	id, err := action.Start(ctx, targetInstance.Task.Definition.ID, "deploy")
	if err != nil {
		return model.Instance{}, err
	}

	go func() {
		const maxJobRuntime = time.Minute * 10

		bgctx := context.Background()

		sess, err := db.Client().StartSession()
		if err != nil {
			log.Println(err)
			return
		}
		defer sess.EndSession(bgctx)

		jobctx, cancel := context.WithTimeout(bgctx, maxJobRuntime)
		defer cancel()

		go func() {
			<-jobctx.Done()

			sess, err := db.Client().StartSession()
			if err != nil {
				log.Println(err)
				return
			}
			defer sess.EndSession(context.Background())

			err = jobctx.Err()

			if err != nil && err != context.Canceled {
				if err := mongo.WithSession(context.Background(), sess, func(sctx mongo.SessionContext) error {
					return action.Stop(sctx, id, fmt.Errorf("deployment timed out after %d minutes", maxJobRuntime))
				}); err != nil {
					fmt.Println("ended?????")
					log.Println(err)
				}
			}
		}()

		if err := mongo.WithSession(jobctx, sess, func(sctx mongo.SessionContext) error {
			//time.Sleep(time.Minute * 5)

			err := deployByType(sctx, id, app.Type, registryInstance, targetInstance, revision, opts.ToBusiness(registryInstance.Repo()))

			return action.Stop(sctx, id, err)

		}); err != nil {
			log.Println(err)
		}
	}()

	instances, err = supplement.Instances(ctx, app.Type, app.GetInstances([]string{target}), false)
	if err != nil {
		return model.Instance{}, err
	}

	cache.Apps.UpdateInstances(app.Name, instances)

	return instances[target], nil
}

func deployByType(ctx mongo.SessionContext, actionID primitive.ObjectID, appType string, source, target model.Instance, revision int64, opts task.DeployOptions) error {

	switch appType {
	case appTypeService:
		return service.Deploy(source, target, revision, opts)
	case appTypeScheduledTask:
		return event.Deploy(source, target, revision, opts)
	case appTypeLambda:
		return lambda.Deploy(source, target, revision, opts)
	case appTypeS3:
		return s3.Deploy(ctx, actionID, source, target, revision, opts)
	default:
		return fmt.Errorf("invalid app type %s", appType)
	}
}
