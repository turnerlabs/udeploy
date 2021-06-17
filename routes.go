package main

import (
	"context"
	"net/http"

	"github.com/turnerlabs/udeploy/component/audit"
	"github.com/turnerlabs/udeploy/component/auth"
	"github.com/turnerlabs/udeploy/component/broker"
	"github.com/turnerlabs/udeploy/component/build"
	"github.com/turnerlabs/udeploy/component/cache"
	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/logger"
	"github.com/turnerlabs/udeploy/component/user"
	"github.com/turnerlabs/udeploy/handler"

	mongosession "github.com/kendavis2/mongo"
	sess "github.com/turnerlabs/udeploy/component/session"

	"github.com/turnerlabs/udeploy/component/db"
	"github.com/turnerlabs/udeploy/component/request"

	"github.com/go-session/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echosession "github.com/turnerlabs/udeploy/component/session"
	"go.mongodb.org/mongo-driver/mongo"
)

const dayInSeconds = 86400

func startRouter(changeNotifier *broker.Broker) {
	e := echo.New()

	//--------------------------------------------------
	//- Middleware
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
	oauth.GET("/logout", handler.Logout)
	oauth.GET("/login", handler.Login)
	oauth.GET("/response", handler.Response)

	//--------------------------------------------------
	//- Redirects
	//--------------------------------------------------
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusMovedPermanently, "/apps") })
	e.GET("/apps/", func(c echo.Context) error { return c.Redirect(http.StatusMovedPermanently, "/apps") })

	//--------------------------------------------------
	//- Broadcast real-time app changes to clients
	//--------------------------------------------------
	events := e.Group("/events", auth.UnAuthError, request.EventStreamContext)
	events.GET("/app/changes", func(c echo.Context) error { return handler.Change(c, changeNotifier) })

	//--------------------------------------------------
	//- V1 routes
	//--------------------------------------------------
	v1 := e.Group("/v1", auth.UnAuthError, request.RouteContext)

	v1.GET("/user", func(c echo.Context) error {
		ctx := c.Get("ctx").(mongo.SessionContext)
		usr := ctx.Value(sess.ContextKey("user")).(user.User)
		return c.JSON(http.StatusOK, usr)
	})
	v1.GET("/users", handler.GetUsers, auth.RequireAdmin)
	v1.POST("/users/:id", handler.SaveUser, auth.RequireAdmin)
	v1.DELETE("/users/:id", handler.DeleteUser, auth.RequireAdmin)

	v1.GET("/notices", handler.GetNotices, auth.RequireAdmin)
	v1.POST("/notices/:id", handler.SaveNotice, auth.RequireAdmin)
	v1.DELETE("/notices/:id", handler.DeleteNotice, auth.RequireAdmin)

	v1.GET("/projects", handler.GetProjects, auth.RequireAdmin)
	v1.POST("/projects/:id", handler.SaveProject, auth.RequireAdmin)
	v1.DELETE("/projects/:id", handler.DeleteProject, auth.RequireAdmin)

	v1.GET("/apps", handler.GetApps)
	v1.POST("/apps/filter", handler.FilterApps)
	v1.GET("/apps/:app", handler.GetApp)
	v1.PUT("/apps/:app/cache", handler.GetApp)
	v1.DELETE("/apps/:app", handler.DeleteApp, auth.RequireAdmin)
	v1.POST("/apps/:app", handler.SaveApp)
	v1.GET("/apps/:app/instances/:registryInstance/registry", build.GetBuilds)

	v1.PUT("/apps/:app/instances/:instance/restart/:desiredCount", handler.Restart, auth.RequireScale, cache.EnsureCache, audit.Audit)
	v1.PUT("/apps/:app/instances/:instance/scale/:desiredCount", handler.Scale, auth.RequireScale, cache.EnsureCache, audit.Audit)

	v1.POST("/apps/:app/instances/:instance/deploy/:registryInstance/:revision", handler.DeployRevision, auth.RequireDeploy, cache.EnsureCache, audit.Audit)

	v1.GET("/apps/:app/instances/:instance/commits", handler.GetInstanceCommits, cache.EnsureCache)
	v1.GET("/apps/:app/instances/:instance/audit", handler.GetAuditEntries)
	v1.GET("/apps/:app/version/range/:current/to/:target/commits", handler.GetVersionCommitsByRange, cache.EnsureCache)

	v1.GET("/apps/:app/instances/:instance/config", handler.GetConfigValue, auth.RequireEdit)
	v1.POST("/apps/:app/instances/:instance/config", handler.SaveConfigValue, auth.RequireEdit)

	//--------------------------------------------------
	//- UI static files
	//--------------------------------------------------
	ui := e.Group("/apps", auth.UnAuthRedirect, cache.NoCache)

	ui.Static("", "vue/pages/apps")
	ui.Static("/:app", "vue/pages/app")
	ui.Static("/:app/instance/:instance", "vue/pages/instance")
	ui.Static("/users", "vue/pages/users")
	ui.Static("/notices", "vue/pages/notices")
	ui.Static("/projects", "vue/pages/projects")

	e.Static("/*", "vue")

	e.Logger.Fatal(e.Start(":8080"))
}
