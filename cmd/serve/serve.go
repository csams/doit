package serve

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/csams/doit/pkg/errors"
	"github.com/csams/doit/pkg/server"
	"github.com/csams/doit/pkg/server/routes"
	"github.com/csams/doit/pkg/storage"
)

func NewCommand() *cobra.Command {
	options := NewOptions()
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the TODO server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(viper.GetViper()); err != nil {
				return err
			}

			if errs := options.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config := NewConfig(options).Complete()

			db, err := storage.New(config.Storage)
			if err != nil {
				return err
			}

			handler := routes.NewHandler(db)
			server, err := server.New(config.Server, handler)
			if err != nil {
				return err
			}
			return server.PrepareRun().Run()
		},
	}

	options.AddFlags(cmd.Flags())
	return cmd
}
