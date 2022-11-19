package server

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	Address string

	CertFile string
	KeyFile  string

	SecureServing bool
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Address, "addr", "a", "0.0.0.0:8080", "the host and port on which to listen")
	fs.StringVarP(&o.CertFile, "cert-file", "c", "", "the file containing the server's serving certificate")
	fs.StringVarP(&o.KeyFile, "key-file", "k", "", "the file containing the server's private key for the serving cert")
}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) Complete(v *viper.Viper) error {
	o.SecureServing = o.CertFile != "" && o.KeyFile != ""
	return nil
}
