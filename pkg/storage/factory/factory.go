package factory

import (
	"fmt"

	"github.com/csams/doit/pkg/storage"
	"github.com/spf13/viper"
)

// New creates a Storage according to the configuration stored in the viper
// instance
func New(v *viper.Viper) (storage.Storage, error) {
	typ := v.GetString("storage.type")
	Factory, exists := registry[typ]
	if !exists {
		return nil, fmt.Errorf("storage type [%s] doesn't exist", typ)
	}
	config := v.Sub("storage.config")
	return Factory(config)
}

// Register is used by the different storage backends to register creation
// functions
func Register(typ string, f FactoryFunc) {
	if _, exists := registry[typ]; exists {
		panic(fmt.Errorf("duplicate backend: [%s]", typ))
	}
	registry[typ] = f
}

type ContextKeyType string

var ContextKey = ContextKeyType("store")

type FactoryFunc func(*viper.Viper) (storage.Storage, error)

var registry = make(map[string]FactoryFunc)
