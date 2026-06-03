package task

import (
	"testing"
	"time"

	"github.com/sho0pi/tickli/internal/types"
)

// applyDateFields covers two deterministic concerns here: RFC3339 start/due
// parsing and the all-day toggle. The natural-language --date path depends on
// time.Now() and is left to manual/integration verification.
func TestApplyDateFields(t *testing.T) {
	t.Run("empty values leave task untouched", func(t *testing.T) {
		task := &types.Task{}
		if err := applyDateFields(task, "", "", "", "", false, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !time.Time(task.StartDate).IsZero() || !time.Time(task.DueDate).IsZero() || task.TimeZone != "" || task.IsAllDay {
			t.Errorf("task should be untouched, got %+v", task)
		}
	})

	t.Run("RFC3339 start and due", func(t *testing.T) {
		task := &types.Task{}
		err := applyDateFields(task, "", "2026-06-02T09:00:00Z", "2026-06-02T17:00:00Z", "Asia/Tokyo", false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := time.Time(task.StartDate).UTC(); got.Hour() != 9 {
			t.Errorf("start hour = %d, want 9", got.Hour())
		}
		if got := time.Time(task.DueDate).UTC(); got.Hour() != 17 {
			t.Errorf("due hour = %d, want 17", got.Hour())
		}
		if task.TimeZone != "Asia/Tokyo" {
			t.Errorf("tz = %q, want Asia/Tokyo", task.TimeZone)
		}
	})

	t.Run("invalid RFC3339 returns error", func(t *testing.T) {
		if err := applyDateFields(&types.Task{}, "", "not-a-time", "", "", false, false); err == nil {
			t.Fatal("expected error for malformed start date")
		}
	})

	t.Run("all-day only applied when changed", func(t *testing.T) {
		task := &types.Task{IsAllDay: true}
		// allDay=false but allDayChanged=false: must stay true.
		if err := applyDateFields(task, "", "", "", "", false, false); err != nil {
			t.Fatal(err)
		}
		if !task.IsAllDay {
			t.Error("IsAllDay should be preserved when flag not changed")
		}
		// allDayChanged=true: value is applied.
		if err := applyDateFields(task, "", "", "", "", false, true); err != nil {
			t.Fatal(err)
		}
		if task.IsAllDay {
			t.Error("IsAllDay should be set to false when flag changed")
		}
	})
}
