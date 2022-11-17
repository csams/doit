package search

import (
	"fmt"

	"github.com/csams/doit/pkg/commands"
	storage "github.com/csams/doit/pkg/storage"
	factory "github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/cobra"
)

func NewCommand(use string) *cobra.Command {
	return &cobra.Command{
		Use: use,
		RunE: func(cmd *cobra.Command, args []string) error {
			st := cmd.Context().Value(factory.ContextKey).(storage.Storage)
			tasks, err := st.Search(&commands.Search{})
			if err != nil {
				return err
			}
			for _, t := range tasks {
				fmt.Printf("%+v\n", t)
			}
			return nil
		},
	}
}
