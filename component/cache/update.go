package cache

import (
	"fmt"

	"github.com/turnerlabs/udeploy/component/app"
)

func (c *appCache) UpdateInstances(appName string, instances map[string]app.Instance) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if cachedApp, appFound := c.apps[appName]; appFound {
		for cachedName, cachedInst := range cachedApp.Instances {
			if instance, instFound := instances[cachedName]; instFound {
				instance.BackupState(cachedInst.CurrentState)
				instances[cachedName] = instance
			}
		}
	}

	for name, inst := range instances {
		c.apps[appName].Instances[name] = inst.Copy()
	}

	c.Notifications <- c.apps[appName].Copy()
}

func (c *appCache) Update(app app.Application) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	for name, i := range app.Instances {
		if len(i.Task.Definition.ID) == 0 {
			return fmt.Errorf("%s %s missing task definition id", app.Name, name)
		}

		c.lookup[i.Task.Definition.ID] = app.Name
	}

	if cachedApp, appFound := c.apps[app.Name]; appFound {
		for cachedName, cachedInst := range cachedApp.Instances {
			if instance, instFound := app.Instances[cachedName]; instFound {
				instance.BackupState(cachedInst.CurrentState)
				app.Instances[cachedName] = instance
			}
		}
	}

	c.apps[app.Name] = app.Copy()

	c.Notifications <- app.Copy()

	return nil
}
