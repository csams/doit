package remove

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/csams/doit/pkg/storage"
	"github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Args: cobra.ExactArgs(1),
		Use:  "remove",
		RunE: removeTask,
	}
}

func removeTask(cmd *cobra.Command, args []string) error {
	store, ok := cmd.Context().Value(factory.ContextKey).(storage.Storage)
	if !ok {
		return errors.New("couldn't retrieve storage object from context")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse task id to remove: %e", err)
	}

	return store.Delete(uint(id))
}
