package cli

import (
	"net/http"

	"github.com/csams/doit/pkg/auth"
	"github.com/go-logr/logr"
	"github.com/rivo/tview"
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
	App    *tview.Application
	Client *http.Client
	Tokens *auth.TokenProvider
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

	if c.App == nil {
		c.App = tview.NewApplication()
	}

	if c.Tokens == nil {
		var err error
		if c.Tokens, err = auth.NewTokenProvider(completeAuth); err != nil {
			return CompletedConfig{}, err
		}
	}

	if c.Client == nil {
		c.Client = auth.CreateClient(c.Options.InsecureClient)
	}

	return CompletedConfig{&completedConfig{
		Options: c.Options,
		Auth:    completeAuth,
		Common:  c.Common,
	}}, nil
}
