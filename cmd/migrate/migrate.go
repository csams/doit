package migrate

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/csams/doit/pkg/errors"
	"github.com/csams/doit/pkg/storage"
)

func NewCommand(log logr.Logger, options *storage.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Create or migrate the database tables.",
		RunE: func(cmd *cobra.Command, args []string) error {

			// TODO: change the viper instance to the correct viper.Sub
			if err := options.Complete(); err != nil {
				return err
			}

			if errs := options.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config := storage.NewConfig(options).Complete()

			db, err := storage.New(config)
			if err != nil {
				return err
			}

			return storage.Migrate(db)
		},
	}

	options.AddFlags(cmd.Flags())
	return cmd
}
