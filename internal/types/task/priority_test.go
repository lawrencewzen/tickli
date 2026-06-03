package task

import (
	"encoding/json"
	"testing"
)

func TestPrioritySet(t *testing.T) {
	tests := []struct {
		in      string
		want    Priority
		wantErr bool
	}{
		{in: "none", want: PriorityNone},
		{in: "low", want: PriorityLow},
		{in: "medium", want: PriorityMedium},
		{in: "high", want: PriorityHigh},
		{in: "HIGH", want: PriorityHigh}, // case-insensitive
		{in: "urgent", wantErr: true},
		{in: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var p Priority
			err := p.Set(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Set(%q) expected error, got %v", tt.in, p)
				}
				return
			}
			if err != nil {
				t.Fatalf("Set(%q) unexpected error: %v", tt.in, err)
			}
			if p != tt.want {
				t.Errorf("Set(%q) = %d, want %d", tt.in, p, tt.want)
			}
		})
	}
}

func TestPriorityUnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   string
		want Priority
	}{
		{in: "0", want: PriorityNone},
		{in: "1", want: PriorityLow},
		{in: "3", want: PriorityMedium},
		{in: "5", want: PriorityHigh},
		{in: "7", want: PriorityNone}, // unknown value clamps to none
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var p Priority
			if err := json.Unmarshal([]byte(tt.in), &p); err != nil {
				t.Fatalf("unmarshal %q: %v", tt.in, err)
			}
			if p != tt.want {
				t.Errorf("unmarshal %q = %d, want %d", tt.in, p, tt.want)
			}
		})
	}
}

func TestPriorityMarshalRoundTrip(t *testing.T) {
	for _, p := range []Priority{PriorityNone, PriorityLow, PriorityMedium, PriorityHigh} {
		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("marshal %d: %v", p, err)
		}
		var got Priority
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal %s: %v", data, err)
		}
		if got != p {
			t.Errorf("round trip: got %d, want %d", got, p)
		}
	}
}
