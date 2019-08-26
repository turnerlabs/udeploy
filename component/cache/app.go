package cache

import (
	"sync"

	"github.com/turnerlabs/udeploy/component/app"
)

// Apps ...
var Apps appCache

func init() {
	Apps = appCache{
		apps:          map[string]app.Application{},
		lookup:        map[string]string{},
		Notifications: make(chan app.Application),
	}
}

type appCache struct {
	apps   map[string]app.Application
	lookup map[string]string
	mux    sync.Mutex

	Notifications chan app.Application
}
