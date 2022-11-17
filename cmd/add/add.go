package add

import (
	"errors"

	util "github.com/csams/doit/cmd/util"
	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/commands"
	storage "github.com/csams/doit/pkg/storage"
	factory "github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Use:  "add",
		RunE: addTask,
	}

	cmd.Flags().StringP("due", "d", "", "due date")
	cmd.Flags().StringP("priority", "p", "", "priority")
	cmd.Flags().StringP("status", "s", "", "status")
	cmd.Flags().StringSlice("tags", nil, "Comma separated list of tags")

	return cmd
}

func addTask(cmd *cobra.Command, args []string) error {
	description := args[0]
	if len(description) == 0 {
		return errors.New("description can't be empty")
	}

	add := &commands.Create{
		Description: description,
	}

	if err := applyAddFlags(add, cmd.Flags()); err != nil {
		return err
	}

	store, ok := cmd.Context().Value(factory.ContextKey).(storage.Storage)
	if !ok {
		return errors.New("couldn't retrieve storage object from context")
	}

	if err := store.Create(add); err != nil {
		return err
	}

	return nil
}

func applyAddFlags(add *commands.Create, flags *pflag.FlagSet) error {
	due, err := util.GetDue(flags)
	if err != nil {
		return err
	}
	add.Due = due

	prio, err := util.GetPriority(flags)
	if err != nil {
		return err
	}
	add.Priority = prio

	status, err := util.GetStatus(flags)
	if err != nil {
		return err
	}
	if status != nil {
		add.Status = *status
	} else {
		add.Status = task.Todo
	}

	tags, err := util.GetTags(flags)
	if err != nil {
		return err
	}
	add.Tags = tags

	return nil
}
