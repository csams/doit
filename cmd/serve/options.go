package serve

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/csams/doit/pkg/server"
	"github.com/csams/doit/pkg/storage"
)

type Options struct {
	Storage *storage.Options
	Server  *server.Options
}

func NewOptions() *Options {
	return &Options{
		Storage: storage.NewOptions(),
		Server:  server.NewOptions(),
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Storage.AddFlags(fs)
	o.Server.AddFlags(fs)
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.Storage.Validate()...)
	errs = append(errs, o.Server.Validate()...)
	return errs
}

func (o *Options) Complete(v *viper.Viper) error {
	// TODO: thread the correct v.Sub into the options where appropriate
	if err := o.Storage.Complete(v); err != nil {
		return err
	}
	if err := o.Server.Complete(v); err != nil {
		return err
	}
	return nil
}
