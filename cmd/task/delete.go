package task

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/completion"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	projectID string
	taskID    string
}

func newDeleteCommand(client *api.Client) *cobra.Command {
	opts := &deleteOptions{}
	cmd := &cobra.Command{
		Use:     "delete <task-id>",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a task",
		Long: `Permanently delete a task by its ID from the current or specified project.

This removes the task immediately without confirmation, so it works in scripts
and non-interactive contexts.`,
		Example: `  # Delete a task in the current project
  tickli task delete abc123def456

  # Delete a task in a specific project
  tickli task delete abc123def456 --project-id xyz789`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.TaskIDs(projectID),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.projectID = projectID
			opts.taskID = args[0]
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteTask(opts.projectID, opts.taskID); err != nil {
				return errors.Wrap(err, "failed to delete task")
			}

			fmt.Printf("%s Task %s deleted\n", color.Red.Sprint("✗"), opts.taskID)
			return nil
		},
	}

	return cmd
}
