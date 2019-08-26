package cache

import (
	"github.com/turnerlabs/udeploy/component/app"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *appCache) GetAll() map[string]app.Application {
	c.mux.Lock()
	defer c.mux.Unlock()

	copies := map[string]app.Application{}

	for name, app := range c.apps {
		copies[name] = app.Copy()
	}

	return copies
}

func (c *appCache) Get(appName string) (app.Application, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	app, found := c.apps[appName]

	return app.Copy(), found
}

func (c *appCache) GetByID(appID primitive.ObjectID) app.Application {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, app := range c.apps {
		if appID == app.ID {
			return app.Copy()
		}
	}

	return app.Application{}
}

func (c *appCache) GetByDefinitionID(taskDefinition string) (app.Application, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	appName, found := c.lookup[taskDefinition]
	if !found {
		return app.Application{}, false
	}

	a, found := c.apps[appName]
	if !found {
		return app.Application{}, false
	}

	return a.Copy(), found
}
