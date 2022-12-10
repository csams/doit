package tui

import (
	"strings"

	"github.com/csams/doit/pkg/auth"
	"github.com/csams/doit/pkg/tui/client"
	"github.com/go-logr/logr"
)

type Config struct {
	Options *Options
	Auth    *auth.Config

	Common
}

type completedConfig struct {
	Options *Options
	Auth    auth.CompletedConfig

	Common
}

type Common struct {
	Client client.Client
	Log    logr.Logger
}

// CompletedConfig can be constructed only from Config.Complete
type CompletedConfig struct {
	*completedConfig
}

func NewConfig(o *Options, log logr.Logger) *Config {
	return &Config{
		Options: o,
		Auth:    auth.NewConfig(o.Auth),
		Common: Common{
			Log: log,
		},
	}
}

func (c *Config) Complete() (CompletedConfig, error) {
	completeAuth := c.Auth.Complete()

	if c.Client.Tokens == nil {
		var err error
		if c.Client.Tokens, err = auth.NewTokenProvider(completeAuth); err != nil {
			return CompletedConfig{}, err
		}
	}

	if c.Client.Http == nil {
		c.Client.Http = auth.NewClient(c.Options.InsecureClient)
	}

	var baseUrl = c.Options.Address
	if !strings.HasSuffix(c.Options.Address, "/") {
		baseUrl = baseUrl + "/"
	}
	c.Client.BaseUrl = baseUrl

	return CompletedConfig{&completedConfig{
		Options: c.Options,
		Auth:    completeAuth,
		Common:  c.Common,
	}}, nil
}
