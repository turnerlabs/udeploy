package main

import (
	"context"
	"log"

	"github.com/turnerlabs/udeploy/component/sync"
	"go.mongodb.org/mongo-driver/mongo"
)

func monitorChanges(ctx context.Context, sess mongo.Session) {

	//--------------------------------------------------
	//- Watch for database app changes
	//--------------------------------------------------
	go func() {
		if err := mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(sync.WatchDatabaseApps(sctx))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Watch for database action changes
	//--------------------------------------------------
	go func() {
		if err := mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(sync.WatchDatabaseActions(sctx))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Watch for task changes
	//--------------------------------------------------
	go func() {
		if err := mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(sync.AWSWatchEvents(sctx))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Watch for cloudwatch alarms
	//--------------------------------------------------
	go func() {
		if err := mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(sync.AWSWatchAlarms(sctx))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()

	//--------------------------------------------------
	//- Watch for s3 changes
	//--------------------------------------------------
	go func() {
		if err := mongo.WithSession(ctx, sess, func(sctx mongo.SessionContext) error {
			log.Fatal(sync.AWSWatchS3(sctx))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()
}
