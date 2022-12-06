package tui

import (
	"github.com/csams/doit/pkg/auth"
	"github.com/spf13/pflag"
)

type Options struct {
	Auth           *auth.Options `mapstructure:"auth"`
	Address        string        `mapstructure:"addr"`
	InsecureClient bool          `mapstructure:"insecure-client"`
}

func NewOptions() *Options {
	return &Options{
		Auth:           auth.NewOptions(),
		Address:        "http://localhost:9090",
		InsecureClient: false,
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.String("client.addr", "http://localhost:9090", "the URL at which the application is hosted")
	fs.Bool("client.insecure-client", false, "the URL at which the application is hosted")

	o.Auth.AddFlags(fs, "client.auth")
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.Auth.Validate()...)
	return errs
}

func (o *Options) Complete() error {
	return o.Auth.Complete()
}
