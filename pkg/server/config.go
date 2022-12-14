package server

import "github.com/csams/doit/pkg/auth"

type Config struct {
	Options *Options
	Auth    *auth.Config
}

type completedConfig struct {
	Options *Options
	Auth    auth.CompletedConfig
}

// CompletedConfig can be constructed only from Config.Complete
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
