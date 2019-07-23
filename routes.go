package main

import (
	"github.com/turnerlabs/udeploy/component/user"
	"context"
	"fmt"
	"net/http"

	"github.com/turnerlabs/udeploy/component/commit"
	sess "github.com/turnerlabs/udeploy/component/session"

	mongosession "github.com/kendavis2/mongo"

	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/component/request"

	"github.com/go-session/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/turnerlabs/udeploy/component/audit"
	"github.com/turnerlabs/udeploy/component/auth"
	"github.com/turnerlabs/udeploy/component/broker"
	"github.com/turnerlabs/udeploy/component/build"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/delete"
	"github.com/turnerlabs/udeploy/component/deploy"
	"github.com/turnerlabs/udeploy/component/get"
	"github.com/turnerlabs/udeploy/component/logger"
	"github.com/turnerlabs/udeploy/component/notify"
	"github.com/turnerlabs/udeploy/component/save"
	"github.com/turnerlabs/udeploy/component/scale"
	echosession "github.com/turnerlabs/udeploy/component/session"
	"go.mongodb.org/mongo-driver/mongo"
)

const dayInSeconds = 86400

func startRouter(changeNotifier *broker.Broker) {
	e := echo.New()

	//--------------------------------------------------
	//- General middleware set up
	//--------------------------------------------------
	e.Use(logger.LogErrors)
	e.Use(middleware.Recover())
	e.Use(echosession.New(
		session.SetStore(mongosession.NewStoreWithClient(context.Background(), db.Client(), cfg.Get["DB_NAME"], "session")),
		session.SetCookieName(cfg.Get["DB_NAME"]),
		session.SetCookieLifeTime(dayInSeconds*30),
		session.SetExpired(dayInSeconds*30),
		session.SetSecure(true),
		session.SetSign([]byte(cfg.Get["OAUTH_SESSION_SIGN"]))))

	//--------------------------------------------------
	//- Ping route
	//--------------------------------------------------
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Version string `json:"version"`
		}{version})
	})

	//--------------------------------------------------
	//- Client configuration route
	//--------------------------------------------------
	e.GET("/config", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			ConsoleLink string `json:"consoleLink"`
		}{cfg.Get["CONSOLE_LINK"]})
	})

	//--------------------------------------------------
	//- OAuth routes
	//--------------------------------------------------
	oauth := e.Group("/oauth2")
	oauth.GET("/logout", func(c echo.Context) error { return auth.Logout(c) })
	oauth.GET("/login", func(c echo.Context) error { return auth.Login(c) })
	oauth.GET("/response", func(c echo.Context) error { return auth.Response(c) })

	//--------------------------------------------------
	//- Redirects
	//--------------------------------------------------
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusMovedPermanently, "/apps") })
	e.GET("/apps/", func(c echo.Context) error { return c.Redirect(http.StatusMovedPermanently, "/apps") })

	//--------------------------------------------------
	//- Publish real-time app changes to clients
	//--------------------------------------------------
	events := e.Group("/events", auth.UnAuthError)
	events.GET("/app/changes", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")

		fmt.Fprint(c.Response().Writer, "event: open\n\n")
		if f, ok := c.Response().Writer.(http.Flusher); ok {
			f.Flush()
		}

		changes := changeNotifier.Subscribe()
		defer changeNotifier.Unsubscribe(changes)

		return notify.Start(c, changes)
	}, request.Context)

	//--------------------------------------------------
	//- v1 routes
	//--------------------------------------------------
	v1 := e.Group("/v1", auth.UnAuthError)

	v1.GET("/user", func(c echo.Context) error {
		ctx := c.Get("ctx").(mongo.SessionContext)
		usr := ctx.Value(sess.ContextKey("user")).(user.User)

		return c.JSON(http.StatusOK, usr)
	}, request.Context)
	v1.GET("/users", func(c echo.Context) error {
		return get.Users(c)
	}, request.Context, auth.RequireAdmin)
	v1.POST("/users/:id", func(c echo.Context) error {
		return save.User(c)
	}, request.Context, auth.RequireAdmin)
	v1.DELETE("/users/:id", func(c echo.Context) error {
		return delete.User(c)
	}, request.Context, auth.RequireAdmin)
	v1.GET("/notices", func(c echo.Context) error {
		return get.Notices(c)
	}, request.Context, auth.RequireAdmin)
	v1.POST("/notices/:id", func(c echo.Context) error {
		return save.Notice(c)
	}, request.Context, auth.RequireAdmin)
	v1.DELETE("/notices/:id", func(c echo.Context) error {
		return delete.Notice(c)
	}, request.Context, auth.RequireAdmin)
	v1.GET("/apps", func(c echo.Context) error {
		return get.Apps(c)
	}, request.Context)
	v1.GET("/apps/:app", func(c echo.Context) error {
		return get.App(c)
	}, request.Context)
	v1.DELETE("/apps/:app", func(c echo.Context) error {
		return delete.App(c)
	}, request.Context, auth.RequireAdmin)
	v1.PUT("/apps/:app/cache", func(c echo.Context) error {
		return cache.App(c)
	}, request.Context)
	v1.POST("/apps/:app", func(c echo.Context) error {
		return save.App(c)
	}, request.Context)
	v1.GET("/apps/:app/instances/:registryInstance/registry", func(c echo.Context) error {
		return build.GetBuilds(c)
	}, request.Context)
	v1.PUT("/apps/:app/instances/:instance/scale/:desiredCount", func(c echo.Context) error {
		return scale.Instance(c, false)
	}, request.Context, auth.RequireScale, cache.EnsureCache, audit.Audit)
	v1.PUT("/apps/:app/instances/:instance/restart/:desiredCount", func(c echo.Context) error {
		return scale.Instance(c, true)
	}, request.Context, auth.RequireScale, cache.EnsureCache, audit.Audit)
	v1.POST("/apps/:app/instances/:instance/deploy/:registryInstance/:revision", func(c echo.Context) error {
		return deploy.Revision(c)
	}, request.Context, auth.RequireDeploy, cache.EnsureCache, audit.Audit)
	v1.GET("/apps/:app/instances/:instance/commits", func(c echo.Context) error {
		return commit.GetInstanceCommits(c)
	}, request.Context, cache.EnsureCache)
	v1.GET("/apps/:app/version/range/:current/to/:target/commits", func(c echo.Context) error {
		return commit.GetVersionCommitsByRange(c)
	}, request.Context, cache.EnsureCache)
	v1.GET("/apps/:app/instances/:instance/audit", func(c echo.Context) error {
		return audit.GetAuditEntries(c)
	}, request.Context)

	//--------------------------------------------------
	//- UI static files
	//--------------------------------------------------
	ui := e.Group("/apps", auth.UnAuthRedirect, cache.NoCache)

	ui.Static("", "vue/pages/apps")
	ui.Static("/:app", "vue/pages/app")
	ui.Static("/:app/instance/:instance", "vue/pages/instance")
	ui.Static("/users", "vue/pages/users")
	ui.Static("/notices", "vue/pages/notices")

	e.Static("/*", "vue")

	e.Logger.Fatal(e.Start(":8080"))
}
