package cli

import (
	"github.com/csams/doit/pkg/auth"
	"github.com/csams/doit/pkg/cli"
	"github.com/csams/doit/pkg/errors"

	"github.com/go-logr/logr"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func NewCommand(log logr.Logger, options *auth.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(); err != nil {
				return err
			}

			if errs := options.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config := auth.NewConfig(options).Complete()
			flow, err := auth.NewTokenProvider(config)
			if err != nil {
				return err
			}

			app := tview.NewApplication()

			prim, err := cli.GetApplication(log, app, flow)
			if err != nil {
				return err
			}

			app.SetRoot(prim, true)
			return app.Run()
		},
	}

	return cmd
}
