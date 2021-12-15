package registry

import (
	"sync"

	"github.com/gofiber/fiber/v2"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
)

type AppStore struct {
	StoreHandler inf.CRUD
	App          *fiber.App
}

type Registry struct {
	apps  map[string]*AppStore
	mutex sync.RWMutex
}

func (r *Registry) Get(name string) *AppStore {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.apps[name]
}

func (r *Registry) Iter(fn func(*fiber.App)) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for k := range r.apps {
		fn(r.apps[k].App)
	}
}

func (r *Registry) Set(name string, appStore *AppStore) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.apps[name] = appStore
}

var (
	regOnce  sync.Once
	registry *Registry
)

func GetRegistry() *Registry {
	regOnce.Do(func() {
		registry = &Registry{
			apps: make(map[string]*AppStore),
		}
	})

	return registry
}
