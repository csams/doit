package server

import (
	"github.com/spf13/pflag"
)

type Options struct {
	Address string `mapstructure:"addr"`

	CertFile string `mapstructure:"cert-file"`
	KeyFile  string `mapstructure:"key-file"`

	SecureServing bool
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.String("serve.addr", "0.0.0.0:8080", "the host and port on which to listen")
	fs.String("serve.cert-file", "", "the file containing the server's serving certificate")
	fs.String("serve.key-file", "", "the file containing the server's private key for the serving cert")
}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) Complete() error {
	o.SecureServing = o.CertFile != "" && o.KeyFile != ""
	return nil
}
