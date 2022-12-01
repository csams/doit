package cli

import (
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Options struct {
	ClientId               string
	TokenFile              string
	LocalAddr              string
	RedirectURL            string
	AuthorizationServerURL string
	InsecureClient         bool
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ClientId, "client-id", "todo-app", "the clientId issued by the authorization server that represents the application")
	fs.StringVarP(&o.TokenFile, "token-file", "t", "$HOME/.config/doit/oidc-token", "the path to the file that holds the id-token")
	fs.StringVar(&o.LocalAddr, "local-addr", "localhost:8080", "the local address that starts the OAuth2 flow")
	fs.StringVar(&o.RedirectURL, "redirect-url", "http://localhost:8080/callback", "the callback URL")
	fs.StringVar(&o.AuthorizationServerURL, "auth-server-url", "https://localhost/realms/todoapp", "the URL to the authorization server")
	fs.BoolVarP(&o.InsecureClient, "insecure-client", "k", false, "validate authorization server certs?")
}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) Complete(v *viper.Viper) error {
	if o.TokenFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		o.TokenFile = path.Join(home, ".config", "doit", "oidc-token")
	} else {
		o.TokenFile = os.ExpandEnv(o.TokenFile)
	}
	return nil
}
