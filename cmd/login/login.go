package login

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/csams/doit/pkg/errors"
	"github.com/csams/doit/pkg/login"
)

func NewCommand(log logr.Logger, options *login.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login using OAuth2.0/OIDC",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(); err != nil {
				return err
			}

			if errs := options.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config := login.NewConfig(options).Complete()
			flow, err := login.NewOIDCFlow(config)
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
