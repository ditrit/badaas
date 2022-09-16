package registry

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/repository"
)

var localRegistry *Registry

// The registry type
type Registry struct {
	UserRepo repository.UserRepository
}

// The factory the the registry type
func FactoryRegistry(typeDatastore Datastore) (*Registry, error) {
	switch typeDatastore {
	case GormDatastore:
		return createGormRegistry()
	default:
		return nil, fmt.Errorf("this type of registry is not implemented")
	}

}

// Replace the global registry instance
func ReplaceGlobals(reg *Registry) {
	localRegistry = reg
}

// Get the global registry instance
func GetRegistry() *Registry {
	if localRegistry == nil {
		panic("registry is nil")
	}
	return localRegistry
}