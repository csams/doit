package auth

import (
	"net/http"
)

type Config struct {
	*Options
	Client *http.Client
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	*completedConfig
}

func NewConfig(o *Options) *Config {
	return &Config{
		Options: o,
	}
}

func (c *Config) Complete() CompletedConfig {
	if c.Client == nil {
		c.Client = CreateClient(c.InsecureClient)
	}
	return CompletedConfig{&completedConfig{
		c,
	}}
}
