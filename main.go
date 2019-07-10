package main // import "github.com/turnerlabs/udeploy"

import (
	"context"
	"log"

	"github.com/turnerlabs/udeploy/component/broker"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/component/deploy"
	"github.com/turnerlabs/udeploy/component/notify"
	"go.mongodb.org/mongo-driver/mongo"
)

var version = "0.0.0-rc"

func main() {

	//--------------------------------------------------
	//- Initialize connections
	//--------------------------------------------------
	ctx := context.Background()

	if err := db.Connect(ctx, cfg.Get["DB_URI"]); err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect(ctx)

	sess, err := db.Client().StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(ctx)

	changeNotifier := broker.NewBroker()
	go changeNotifier.Start()

	go func() {
		for msg := range cache.Apps.Notifications {
			changeNotifier.Publish(msg)
		}
	}()

	//--------------------------------------------------
	//- Cache applications
	//--------------------------------------------------
	go func() {
		if err = mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			return cache.Ensure(sctx)
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Monitor for changes
	//--------------------------------------------------
	monitorChanges(ctx, sess)

	//--------------------------------------------------
	//- Propagate deployments
	//--------------------------------------------------
	go func() {
		changes := changeNotifier.Subscribe()
		defer changeNotifier.Unsubscribe(changes)

		if err = mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(deploy.Propagate(sctx, changes))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Send SQS queue notifications
	//--------------------------------------------------
	go func() {
		changes := changeNotifier.Subscribe()
		defer changeNotifier.Unsubscribe(changes)

		if err = mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(notify.Watch(sctx, changes))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Routing
	//--------------------------------------------------
	startRouter(changeNotifier)
}
