package storage

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	DSN string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) Complete(v *viper.Viper) error {
	o.DSN = v.GetString("dsn")
	return nil
}
