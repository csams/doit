package serve

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/csams/doit/pkg/auth"
	"github.com/csams/doit/pkg/errors"
	"github.com/csams/doit/pkg/server"
	"github.com/csams/doit/pkg/server/routes"
	"github.com/csams/doit/pkg/storage"
)

func NewCommand(log logr.Logger, storageOptions *storage.Options, serverOptions *server.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the TODO server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := storageOptions.Complete(); err != nil {
				return err
			}

			if err := serverOptions.Complete(); err != nil {
				return err
			}

			if errs := storageOptions.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			if errs := serverOptions.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			storageConfig := storage.NewConfig(storageOptions).Complete()
			serverConfig := server.NewConfig(serverOptions).Complete()

			db, err := storage.New(storageConfig)
			if err != nil {
				return err
			}

			authProvider, err := auth.NewTokenProvider(serverConfig.Auth)
			if err != nil {
				return err
			}

			handler := routes.NewHandler(db, authProvider, log.WithName("rootHandler"))
			server, err := server.New(serverConfig, handler, log.WithName("server"))
			if err != nil {
				return err
			}

			return server.PrepareRun().Run()
		},
	}

	storageOptions.AddFlags(cmd.Flags())
	serverOptions.AddFlags(cmd.Flags())

	return cmd
}
