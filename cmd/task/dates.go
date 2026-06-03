package task

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sho0pi/tickli/internal/types"
	"github.com/sho0pi/tickli/internal/utils"
)

// applyDateFields writes the time-related fields onto t from raw flag values.
//
// date is a natural-language expression ("tomorrow 5pm") that sets start, due
// and all-day together. start/due are RFC3339 timestamps that override the
// corresponding field. tz sets the timezone. allDayChanged reports whether the
// --all-day flag was explicitly set (so its value is only applied on demand).
//
// Empty string values are treated as "not provided" and leave t untouched,
// which is what makes this safe for both create (zero-value task) and update
// (existing task being partially modified).
func applyDateFields(t *types.Task, date, start, due, tz string, allDay, allDayChanged bool) error {
	if date != "" {
		r, err := utils.ParseTimeExpression(date)
		if err != nil {
			return errors.Wrap(err, "failed to parse date range")
		}
		t.StartDate = types.TickTickTime(r.Start())
		t.DueDate = types.TickTickTime(r.End())
		t.IsAllDay = r.IsAllDay()
	}
	if start != "" {
		startDate, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return errors.Wrap(err, "failed to parse start date")
		}
		t.StartDate = types.TickTickTime(startDate)
	}
	if due != "" {
		dueDate, err := time.Parse(time.RFC3339, due)
		if err != nil {
			return errors.Wrap(err, "failed to parse due date")
		}
		t.DueDate = types.TickTickTime(dueDate)
	}
	if tz != "" {
		t.TimeZone = tz
	}
	if allDayChanged {
		t.IsAllDay = allDay
	}
	return nil
}
