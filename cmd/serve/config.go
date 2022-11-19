package serve

import (
	"github.com/csams/doit/pkg/server"
	"github.com/csams/doit/pkg/storage"
)

type Config struct {
	Server  *server.Config
	Storage *storage.Config
}

type completedConfig struct {
	Server  server.CompletedConfig
	Storage storage.CompletedConfig
}

type CompletedConfig struct {
	*completedConfig
}

func NewConfig(o *Options) *Config {
	return &Config{
		Server:  server.NewConfig(o.Server),
		Storage: storage.NewConfig(o.Storage),
	}
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{
		Server:  c.Server.Complete(),
		Storage: c.Storage.Complete(),
	}}
}
