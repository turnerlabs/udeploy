package cache

import "go.mongodb.org/mongo-driver/bson/primitive"
 
func (c *appCache) Remove(appName string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.apps, appName)
}

func (c *appCache) RemoveByID(appID primitive.ObjectID) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, app := range c.apps {
		if appID == app.ID {
			delete(c.apps, app.Name)
		}
	}
}