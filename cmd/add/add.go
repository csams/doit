package add

import (
	"github.com/spf13/cobra"
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
	return nil
}
