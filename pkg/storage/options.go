package storage

import (
	"github.com/spf13/pflag"
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

func (o *Options) Complete() error {
	return nil
}
