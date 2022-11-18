package start

import (
	"errors"
	"strconv"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/commands"
	"github.com/csams/doit/pkg/storage"
	"github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Args: cobra.ExactArgs(1),
		Use:  "start",
		RunE: start,
	}
}

func start(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	started := apis.Status(apis.Doing)
	mod := &commands.Modify{
		Id:     uint(id),
		Status: &started,
	}

	store, ok := cmd.Context().Value(factory.ContextKey).(storage.Storage)
	if !ok {
		return errors.New("couldn't retrieve storage object from context")
	}

	if err := store.Update(mod); err != nil {
		return err
	}

	return nil
}
