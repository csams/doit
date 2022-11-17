package modify

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/csams/doit/cmd/util"
	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/commands"
	storage "github.com/csams/doit/pkg/storage"
	factory "github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "modify",
		Args: cobra.ExactArgs(2),
		RunE: modifyTask,
	}

	cmd.Flags().StringP("due", "d", "", "due date")
	cmd.Flags().StringP("priority", "p", "", "priority")
	cmd.Flags().StringP("status", "s", "", "status")
	cmd.Flags().StringSlice("tags", nil, "Comma separated list of tags")

	return cmd
}

func modifyTask(cmd *cobra.Command, args []string) error {
	update := &commands.Modify{}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("not a valid task id: [%s]; %e", args[0], err)
	}

	update.Id = task.Identity(id)

	description := args[1]
	if description == "" {
		update.Description = nil
	} else {
		update.Description = &description
	}

	if err := applyFlags(update, cmd.Flags()); err != nil {
		return err
	}

	store, ok := cmd.Context().Value(factory.ContextKey).(storage.Storage)
	if !ok {
		return errors.New("couldn't retrieve storage object from context")
	}
	if err := store.Update(update); err != nil {
		return err
	}

	return nil
}

func applyFlags(update *commands.Modify, flags *pflag.FlagSet) error {
	due, err := util.GetDue(flags)
	if err != nil {
		return err
	}

	if due != nil {
		update.Due = due
	}

	prio, err := util.GetPriority(flags)
	if err != nil {
		return err
	}

	if prio != task.Undefined {
		update.Priority = prio
	}

	status, err := util.GetStatus(flags)
	if err != nil {
		return err
	}

	if status != nil {
		update.Status = status
	}

	tags, err := util.GetTags(flags)
	if err != nil {
		return err
	}

	if tags != nil {
		update.Tags = tags
	}

	return nil
}
