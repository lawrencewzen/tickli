package task

import (
	"testing"

	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/types/task"
)

func sampleTasks() []types.Task {
	return []types.Task{
		{ID: "a", Title: "low untagged", Priority: task.PriorityLow},
		{ID: "b", Title: "high work", Priority: task.PriorityHigh, Tags: []string{"work"}},
		{ID: "c", Title: "medium home", Priority: task.PriorityMedium, Tags: []string{"home"}},
		{ID: "d", Title: "none work", Priority: task.PriorityNone, Tags: []string{"work", "urgent"}},
	}
}

func ids(tasks []types.Task) []string {
	out := make([]string, len(tasks))
	for i := range tasks {
		out[i] = tasks[i].ID
	}
	return out
}

func TestFilterTasks(t *testing.T) {
	tests := []struct {
		name    string
		opts    *listOptions
		wantIDs []string
	}{
		{
			name:    "no filters returns all",
			opts:    &listOptions{},
			wantIDs: []string{"a", "b", "c", "d"},
		},
		{
			name:    "priority threshold is inclusive lower bound",
			opts:    &listOptions{priority: task.PriorityMedium},
			wantIDs: []string{"b", "c"}, // high(5) and medium(3) >= medium(3)
		},
		{
			name:    "highest priority only",
			opts:    &listOptions{priority: task.PriorityHigh},
			wantIDs: []string{"b"},
		},
		{
			name:    "tag filter",
			opts:    &listOptions{tag: "work"},
			wantIDs: []string{"b", "d"},
		},
		{
			name:    "priority and tag combined",
			opts:    &listOptions{priority: task.PriorityHigh, tag: "work"},
			wantIDs: []string{"b"},
		},
		{
			name:    "tag with no matches",
			opts:    &listOptions{tag: "nonexistent"},
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ids(filterTasks(sampleTasks(), tt.opts))
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("got %v, want %v", got, tt.wantIDs)
			}
			for i := range tt.wantIDs {
				if got[i] != tt.wantIDs[i] {
					t.Fatalf("got %v, want %v", got, tt.wantIDs)
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	in := []types.Task{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	got := Filter(in, func(tk types.Task) bool { return tk.ID != "2" })
	if len(got) != 2 || got[0].ID != "1" || got[1].ID != "3" {
		t.Errorf("Filter() = %v, want ids 1,3", ids(got))
	}
}
