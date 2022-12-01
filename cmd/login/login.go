package login

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/csams/doit/pkg/cli"
	"github.com/csams/doit/pkg/errors"
)

func NewCommand(log logr.Logger) *cobra.Command {
	options := cli.NewOptions()
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login using OAuth2.0/OIDC",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(viper.GetViper()); err != nil {
				return err
			}

			if errs := options.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config := cli.NewConfig(options).Complete()
			flow, err := cli.NewOIDCFlow(config)
			if err != nil {
				return err
			}
			tok, err := flow.GetIdToken()
			fmt.Println(tok)
			return err
		},
	}

	options.AddFlags(cmd.Flags())
	return cmd
}
