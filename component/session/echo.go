package session

import (
	"github.com/go-session/session"
	"github.com/labstack/echo/v4"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c echo.Context) bool

	// Config defines the config for Session middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper
		// StoreKey keys stored in the context
		StoreKey string
		// ManageKey keys stored in the context
		ManageKey string
	}
)

var (
	// DefaultConfig is the default Recover middleware config.
	DefaultConfig = Config{
		Skipper:   func(_ echo.Context) bool { return false },
		StoreKey:  "github.com/go-session/echo-session/store",
		ManageKey: "github.com/go-session/echo-session/manage",
	}

	storeKey  string
	manageKey string
)

// New create a session middleware
func New(opt ...session.Option) echo.MiddlewareFunc {
	return NewWithConfig(DefaultConfig, opt...)
}

// NewWithConfig create a session middleware
func NewWithConfig(config Config, opt ...session.Option) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	manageKey = config.ManageKey
	if manageKey == "" {
		manageKey = DefaultConfig.ManageKey
	}

	storeKey = config.StoreKey
	if storeKey == "" {
		storeKey = DefaultConfig.StoreKey
	}

	manage := session.NewManager(opt...)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			c.Set(manageKey, manage)
			store, err := manage.Start(nil, c.Response(), c.Request())
			if err != nil {
				return err
			}
			c.Set(storeKey, store)
			return next(c)
		}
	}
}

// FromContext get session storage from context
func FromContext(ctx echo.Context) session.Store {
	return ctx.Get(storeKey).(session.Store)
}

// Destroy a session
func Destroy(ctx echo.Context) error {
	return ctx.Get(manageKey).(*session.Manager).Destroy(nil, ctx.Response(), ctx.Request())
}

// Refresh a session and return to session storage
func Refresh(ctx echo.Context) (session.Store, error) {
	return ctx.Get(manageKey).(*session.Manager).Refresh(nil, ctx.Response(), ctx.Request())
}
