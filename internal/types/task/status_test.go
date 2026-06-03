package task

import (
	"encoding/json"
	"testing"
)

func TestStatusUnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   string
		want Status
	}{
		{in: "0", want: StatusNormal},
		{in: "2", want: StatusComplete},
		{in: "9", want: StatusNormal}, // unknown clamps to normal
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var s Status
			if err := json.Unmarshal([]byte(tt.in), &s); err != nil {
				t.Fatalf("unmarshal %q: %v", tt.in, err)
			}
			if s != tt.want {
				t.Errorf("unmarshal %q = %d, want %d", tt.in, s, tt.want)
			}
		})
	}
}

func TestStatusMarshalRoundTrip(t *testing.T) {
	for _, s := range []Status{StatusNormal, StatusComplete} {
		data, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("marshal %d: %v", s, err)
		}
		var got Status
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal %s: %v", data, err)
		}
		if got != s {
			t.Errorf("round trip: got %d, want %d", got, s)
		}
	}
}
