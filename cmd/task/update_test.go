package task

import (
	"testing"
	"time"

	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/types/task"
)

// changedSet returns a Flags().Changed-style predicate that reports true only
// for the named flags.
func changedSet(names ...string) func(string) bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return func(name string) bool { return m[name] }
}

func existingTask() *types.Task {
	return &types.Task{
		ID:        "tid",
		ProjectID: "pid",
		Title:     "original",
		Content:   "original content",
		Priority:  task.PriorityLow,
		Tags:      []string{"old"},
	}
}

func TestApplyTaskUpdates(t *testing.T) {
	t.Run("only changed fields are applied, rest preserved", func(t *testing.T) {
		tk := existingTask()
		opts := &updateOptions{title: "renamed"}
		if err := applyTaskUpdates(tk, opts, changedSet("title")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tk.Title != "renamed" {
			t.Errorf("Title = %q, want renamed", tk.Title)
		}
		// Everything else must be preserved.
		if tk.Content != "original content" || tk.Priority != task.PriorityLow || len(tk.Tags) != 1 || tk.Tags[0] != "old" {
			t.Errorf("non-title fields changed unexpectedly: %+v", tk)
		}
		// id/projectId must never be touched (projectId empty => duplicate task bug).
		if tk.ID != "tid" || tk.ProjectID != "pid" {
			t.Errorf("identity fields changed: id=%q projectId=%q", tk.ID, tk.ProjectID)
		}
	})

	t.Run("priority and tags", func(t *testing.T) {
		tk := existingTask()
		opts := &updateOptions{priority: task.PriorityHigh, tags: []string{"a", "b"}}
		if err := applyTaskUpdates(tk, opts, changedSet("priority", "tags")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tk.Priority != task.PriorityHigh {
			t.Errorf("Priority = %d, want high", tk.Priority)
		}
		if len(tk.Tags) != 2 || tk.Tags[0] != "a" || tk.Tags[1] != "b" {
			t.Errorf("Tags = %v, want [a b]", tk.Tags)
		}
	})

	t.Run("ISO due date", func(t *testing.T) {
		tk := existingTask()
		opts := &updateOptions{dueDate: "2026-06-02T17:00:00Z"}
		if err := applyTaskUpdates(tk, opts, changedSet("due")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := time.Time(tk.DueDate).UTC(); got.Hour() != 17 {
			t.Errorf("due hour = %d, want 17", got.Hour())
		}
	})

	t.Run("invalid ISO due returns error", func(t *testing.T) {
		tk := existingTask()
		opts := &updateOptions{dueDate: "nope"}
		if err := applyTaskUpdates(tk, opts, changedSet("due")); err == nil {
			t.Fatal("expected error for malformed due date")
		}
	})

	t.Run("nothing changed leaves task identical", func(t *testing.T) {
		tk := existingTask()
		before := *tk
		if err := applyTaskUpdates(tk, &updateOptions{title: "ignored"}, changedSet()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tk.Title != before.Title || tk.Content != before.Content || tk.Priority != before.Priority {
			t.Errorf("task changed despite no flags set: %+v", tk)
		}
	})
}
