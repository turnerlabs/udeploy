package request

import (
	"context"
	"time"

	"github.com/turnerlabs/udeploy/component/auth"

	sess "github.com/turnerlabs/udeploy/component/session"
	"github.com/turnerlabs/udeploy/component/user"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/db"
	"go.mongodb.org/mongo-driver/mongo"
)

const timeoutSeconds = 120

// RouteContext provides user session for a limited timeframe.
func RouteContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctxParent, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
		defer cancel()

		session, err := db.Client().StartSession()
		if err != nil {
			return err
		}

		var ctx context.Context

		if err = mongo.WithSession(ctxParent, session, func(sctx mongo.SessionContext) error {

			usr, err := user.Get(sctx, c.Get(auth.UserIDParam).(string))
			if err != nil {
				return err
			}

			usr, err = user.Inherit(sctx, usr)

			ctx = context.WithValue(ctxParent, sess.ContextKey("user"), usr)

			return nil
		}); err != nil {
			return err
		}

		err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) error {
			c.Set("ctx", sctx)

			return next(c)
		})

		session.EndSession(ctx)

		return err
	}
}

// EventStreamContext provides user session for an indefinite timeframe
func EventStreamContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctxParent, cancel := context.WithCancel(context.Background())
		defer cancel()

		session, err := db.Client().StartSession()
		if err != nil {
			return err
		}

		var ctx context.Context

		if err = mongo.WithSession(ctxParent, session, func(sctx mongo.SessionContext) error {

			usr, err := user.Get(sctx, c.Get(auth.UserIDParam).(string))
			if err != nil {
				return err
			}

			usr, err = user.Inherit(sctx, usr)

			ctx = context.WithValue(ctxParent, sess.ContextKey("user"), usr)

			return nil
		}); err != nil {
			return err
		}

		err = mongo.WithSession(ctx, session, func(sctx mongo.SessionContext) error {
			c.Set("ctx", sctx)

			return next(c)
		})

		session.EndSession(ctx)

		return err
	}
}
