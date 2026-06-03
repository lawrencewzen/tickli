package task

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/types/task"
	"github.com/spf13/cobra"
)

type listOptions struct {
	priority  task.Priority
	tag       string
	projectID string
	output    types.OutputFormat
}

func filterTasks(tasks []types.Task, opts *listOptions) []types.Task {
	// Filter by priority
	tasks = Filter(tasks, func(t types.Task) bool {
		return t.Priority >= opts.priority
	})

	// Filter by tags
	tasks = Filter(tasks, func(t types.Task) bool {
		if opts.tag != "" {
			return slices.Contains(t.Tags, opts.tag)
		}
		return true
	})

	return tasks
}

func newListCommand(client *api.Client) *cobra.Command {
	opts := &listOptions{output: types.OutputSimple}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List tasks in a project",
		Long: `Print tasks in the current project or a specified project, one per line.

You can filter tasks by priority and tag. The output is plain text (or JSON
with -o json) so it can be piped into other commands or scripts.`,
		Example: `  # List tasks in current project
  tickli task list

  # List tasks with specific tag
  tickli task list -t important

  # List high priority tasks
  tickli task list -p high

  # List tasks in specific project
  tickli task list --project-id abc123def456`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks, err := client.ListTasks(projectID)
			if err != nil {
				return errors.Wrap(err, "failed to fetch tasks")
			}

			filteredTasks := filterTasks(tasks, opts)

			if opts.output == types.OutputJSON {
				if filteredTasks == nil {
					filteredTasks = []types.Task{}
				}
				data, err := json.MarshalIndent(filteredTasks, "", "  ")
				if err != nil {
					return errors.Wrap(err, "failed to marshal tasks")
				}
				fmt.Println(string(data))
				return nil
			}

			if len(filteredTasks) == 0 {
				fmt.Println("No tasks found.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tPRIORITY\tSTATUS\tTITLE")
			for i := range filteredTasks {
				t := filteredTasks[i]
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.ID, t.Priority.String(), t.Status.String(), t.Title)
			}
			return w.Flush()
		},
	}
	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "Only show tasks with this specific tag")
	cmd.Flags().VarP(&opts.priority, "priority", "p", "Only show tasks with this priority level or higher")
	_ = cmd.RegisterFlagCompletionFunc("priority", task.PriorityCompletionFunc)
	cmd.Flags().VarP(&opts.output, "output", "o", "Display format: simple (human-readable) or json (machine-readable)")
	_ = cmd.RegisterFlagCompletionFunc("output", types.OutputFormatCompletionFunc)

	return cmd
}

func Filter(tasks []types.Task, predicate func(task types.Task) bool) []types.Task {
	var result []types.Task
	for _, t := range tasks {
		if predicate(t) {
			result = append(result, t)
		}
	}
	return result
}
