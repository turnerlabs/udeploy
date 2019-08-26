package deploy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/turnerlabs/udeploy/component/app"

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
)

const (
	appTypeService       = "service"
	appTypeScheduledTask = "scheduled-task"
	appTypeLambda        = "lambda"
	appTypeS3            = "s3"

	maxJobRuntime = time.Minute * 10
)

// Options ...
type Options struct {
	Override bool              `json:"override"`
	Env      map[string]string `json:"env"`
	Secrets  map[string]string `json:"secrets"`

	ImageTag string `json:"imageTag"`
}

// ToBusiness ...
func (o Options) ToBusiness(repository string) task.DeployOptions {
	m := task.DeployOptions{
		Environment: o.Env,
		Secrets:     o.Secrets,
		Override:    o.Override,
	}

	if len(repository) > 0 && len(o.ImageTag) > 0 {
		m.Image = fmt.Sprintf("%s:%s", repository, o.ImageTag)
	}

	return m
}

// Deploy ...
func Deploy(ctx mongo.SessionContext, application app.Application, target, source string, revision int64, opts Options) (app.Instance, error) {

	instances := application.GetInstances([]string{target, source})

	instances, err := supplement.Instances(ctx, application.Type, instances, false)
	if err != nil {
		return app.Instance{}, err
	}

	targetInstance, targetExists := instances[target]
	if !targetExists {
		return app.Instance{}, err
	}

	registryInstance, sourceExists := instances[source]
	if !sourceExists {
		return app.Instance{}, err
	}

	id, err := action.Start(ctx, targetInstance.Task.Definition.ID, "deploy", maxJobRuntime)
	if err != nil {
		return app.Instance{}, err
	}

	go func() {
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
					log.Println(err)
				}
			}
		}()

		if err := mongo.WithSession(jobctx, sess, func(sctx mongo.SessionContext) error {
			//time.Sleep(time.Minute * 5)

			err := deployByType(sctx, id, application.Type, registryInstance, targetInstance, revision, opts.ToBusiness(registryInstance.Repo()))

			return action.Stop(sctx, id, err)

		}); err != nil {
			log.Println(err)
		}
	}()

	instances, err = supplement.Instances(ctx, application.Type, application.GetInstances([]string{target}), false)
	if err != nil {
		return app.Instance{}, err
	}

	cache.Apps.UpdateInstances(application.Name, instances)

	return instances[target], nil
}

func deployByType(ctx mongo.SessionContext, actionID primitive.ObjectID, appType string, source, target app.Instance, revision int64, opts task.DeployOptions) error {

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
