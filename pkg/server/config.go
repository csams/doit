package server

import "github.com/csams/doit/pkg/auth"

type Config struct {
	*Options
	Auth *auth.Config
}

type completedConfig struct {
	*Options
	Auth auth.CompletedConfig
}

type CompletedConfig struct {
	*completedConfig
}

func NewConfig(o *Options) *Config {
	return &Config{
		Options: o,
		Auth:    auth.NewConfig(o.Auth),
	}
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{
		Options: c.Options,
		Auth:    c.Auth.Complete(),
	}}
}
