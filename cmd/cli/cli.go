package cli

import (
	"github.com/csams/doit/pkg/cli"
	"github.com/csams/doit/pkg/errors"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

func NewCommand(log logr.Logger, options *cli.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(); err != nil {
				return err
			}

			if errs := options.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config, err := cli.NewConfig(options, log).Complete()
			if err != nil {
				return err
			}

			c, err := cli.New(config)
			if err != nil {
				return err
			}

			c.App.SetRoot(c.Root, true)
			return c.App.Run()
		},
	}

	options.AddFlags(cmd.Flags())

	return cmd
}
