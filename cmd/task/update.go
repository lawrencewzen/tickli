package task

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sho0pi/tickli/internal/api"
	"github.com/sho0pi/tickli/internal/completion"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/types/task"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	title    string
	content  string
	priority task.Priority
	tags     []string

	allDay    bool
	date      string
	startDate string
	dueDate   string
	timeZone  string

	projectID string
	taskID    string
}

// fieldFlags are the flags that carry an editable task field. At least one must
// be set for an update to make sense.
var fieldFlags = []string{"title", "content", "priority", "tags", "date", "start", "due", "tz", "all-day"}

// applyTaskUpdates overwrites only the fields whose flag the user changed, onto
// an existing task fetched from the API. It never touches t.ID or t.ProjectID:
// those come from the fetched task and must reach the update endpoint intact,
// otherwise TickTick creates a duplicate task instead of updating.
func applyTaskUpdates(t *types.Task, opts *updateOptions, changed func(name string) bool) error {
	if changed("title") {
		t.Title = opts.title
	}
	if changed("content") {
		t.Content = opts.content
	}
	if changed("priority") {
		t.Priority = opts.priority
	}
	if changed("tags") {
		t.Tags = opts.tags
	}
	return applyDateFields(t, opts.date, opts.startDate, opts.dueDate, opts.timeZone, opts.allDay, changed("all-day"))
}

func newUpdateCommand(client *api.Client) *cobra.Command {
	opts := &updateOptions{}
	cmd := &cobra.Command{
		Use:     "update <task-id>",
		Aliases: []string{"edit", "set"},
		Short:   "Update an existing task",
		Long: `Update fields of an existing task in the current or specified project.

Only the fields you pass are changed; everything else is preserved. The task is
fetched first and re-submitted with your changes, so its project association is
kept intact. At least one field flag is required.`,
		Example: `  # Rename a task
  tickli task update abc123 -t "New title"

  # Bump priority and add tags
  tickli task update abc123 -p high --tags urgent,work

  # Set a due date
  tickli task update abc123 --due "2026-02-18T18:00:00Z"`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.TaskIDs(projectID),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.projectID = projectID
			opts.taskID = args[0]
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			changed := false
			for _, f := range fieldFlags {
				if cmd.Flags().Changed(f) {
					changed = true
					break
				}
			}
			if !changed {
				return errors.New("nothing to update; pass at least one field flag (e.g. -t, -p, --tags)")
			}

			existing, err := client.GetTask(opts.projectID, opts.taskID)
			if err != nil {
				return errors.Wrap(err, "failed to fetch task")
			}
			if existing.ID != opts.taskID {
				return fmt.Errorf("task %s not found for project %s", opts.taskID, opts.projectID)
			}

			if err := applyTaskUpdates(existing, opts, cmd.Flags().Changed); err != nil {
				return err
			}
			// Guard against the duplicate-task bug if the API ever returns a task
			// without its project id.
			if existing.ProjectID == "" {
				existing.ProjectID = opts.projectID
			}

			updated, err := client.UpdateTask(existing)
			if err != nil {
				return errors.Wrap(err, "failed to update task")
			}

			fmt.Printf("Updated task %s\n", updated.ID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.title, "title", "t", "", "New title")
	cmd.Flags().StringVarP(&opts.content, "content", "c", "", "New content")
	cmd.Flags().VarP(&opts.priority, "priority", "p", "New priority: none, low, medium, high")
	_ = cmd.RegisterFlagCompletionFunc("priority", task.PriorityCompletionFunc)
	cmd.Flags().StringSliceVar(&opts.tags, "tags", []string{}, "Replace tags (comma-separated)")
	cmd.Flags().BoolVarP(&opts.allDay, "all-day", "a", false, "Set as an all-day task without specific time")
	cmd.Flags().StringVar(&opts.date, "date", "", "Set date with natural language (e.g., 'today', 'next week')")
	cmd.Flags().StringVar(&opts.startDate, "start", "", "When the task begins (ISO format: '2025-02-18T15:04:05Z')")
	cmd.Flags().StringVar(&opts.dueDate, "due", "", "When the task is due (ISO format: '2025-02-18T18:00:00Z')")
	cmd.Flags().StringVar(&opts.timeZone, "tz", "", "Timezone for date calculations (e.g., 'America/Los_Angeles')")

	cmd.MarkFlagsMutuallyExclusive("date", "all-day")
	cmd.MarkFlagsMutuallyExclusive("date", "start")
	cmd.MarkFlagsMutuallyExclusive("date", "due")

	return cmd
}
