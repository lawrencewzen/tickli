package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTickTickTimeUnmarshalJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var tt TickTickTime
		if err := json.Unmarshal([]byte(`"2026-06-02T15:04:05+0000"`), &tt); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := time.Time(tt).UTC()
		if got.Year() != 2026 || got.Month() != time.June || got.Day() != 2 || got.Hour() != 15 {
			t.Errorf("parsed time = %v, want 2026-06-02 15:04:05 UTC", got)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		var tt TickTickTime
		if err := json.Unmarshal([]byte(`"not-a-time"`), &tt); err == nil {
			t.Fatal("expected error for malformed time string")
		}
	})
}

func TestTickTickTimeMarshalRoundTrip(t *testing.T) {
	const raw = `"2026-06-02T15:04:05+0000"`
	var tt TickTickTime
	if err := json.Unmarshal([]byte(raw), &tt); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	data, err := json.Marshal(tt)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != raw {
		t.Errorf("round trip = %s, want %s", data, raw)
	}
}
