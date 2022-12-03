package auth

import (
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/pflag"
)

type Options struct {
	ClientId               string `mapstructure:"client-id"`
	TokenFile              string `mapstructure:"token-file"`
	LocalAddr              string `mapstructure:"local-addr"`
	RedirectURL            string `mapstructure:"redirect-url"`
	AuthorizationServerURL string `mapstructure:"server-url"`
	InsecureClient         bool   `mapstructure:"insecure-client"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	fs.String(prefix+"client-id", "todo-app", "the clientId issued by the authorization server that represents the application")
	fs.StringP(prefix+"token-file", "t", "$HOME/.config/doit/oidc-token", "the path to the file that holds the id-token")
	fs.String(prefix+"local-addr", "localhost:8080", "the local address that starts the OAuth2 flow")
	fs.String(prefix+"redirect-url", "http://localhost:8080/callback", "the callback URL")
	fs.String(prefix+"server-url", "https://localhost/realms/todoapp", "the URL to the authorization server")
	fs.BoolP(prefix+"insecure-client", "k", false, "validate authorization server certs?")
}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) Complete() error {
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
