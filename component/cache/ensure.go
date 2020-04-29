package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/app"
	"github.com/turnerlabs/udeploy/component/integration/aws/secretsmanager"
	"github.com/turnerlabs/udeploy/component/supplement"
	"go.mongodb.org/mongo-driver/mongo"
)

// EnsureCache ...
func EnsureCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Get("ctx").(mongo.SessionContext)

		if err := EnsureApp(ctx, c.Param("app")); err != nil {
			return err
		}

		return next(c)
	}
}

// EnsureApp ...
func EnsureApp(ctx mongo.SessionContext, appName string) error {
	if _, found := Apps.Get(appName); !found {
		apps, err := app.Get(ctx, appName)
		if err != nil {
			return err
		}

		if len(apps) == 0 {
			return fmt.Errorf("%s not found", appName)
		}

		dbApp := apps[0]

		instances, err := supplement.Instances(ctx, dbApp.Type, dbApp.Instances, false)
		if err != nil {
			return err
		}

		dbApp.Instances = instances

		token, err := secretsmanager.Get(apps[0].RepoAccessTokenKey())
		if err != nil {
			return err
		}

		dbApp.Repo.AccessToken = token

		return Apps.Update(dbApp)
	}

	return nil
}

// Ensure ...
func Ensure(ctx mongo.SessionContext) error {
	apps, err := app.Get(ctx, "")
	if err != nil {
		return err
	}

	total := len(apps)
	processed := 0
	cached := 0

	log.Printf("APP_CACHE: caching %d apps\n", total)

	for _, a := range apps {
		if err := EnsureApp(ctx, a.Name); err != nil {
			log.Printf("APP_CACHE: %s\n", err)
		} else {
			cached++
		}

		processed++

		if processed%10 == 0 || processed == total {
			log.Printf("APP_CACHE: %d cached\n", cached)
		}

		// Wait to avoid hitting AWS api rate limits.
		// There is no hurry caching apps since they are also
		// loaded into cache when a user views the portal.
		time.Sleep(1 * time.Second)
	}

	return nil
}
