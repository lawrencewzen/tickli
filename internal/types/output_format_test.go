package types

import "testing"

func TestOutputFormatSet(t *testing.T) {
	tests := []struct {
		in      string
		want    OutputFormat
		wantErr bool
	}{
		{in: "simple", want: OutputSimple},
		{in: "json", want: OutputJSON},
		{in: "yaml", wantErr: true},
		{in: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var o OutputFormat
			err := o.Set(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Set(%q) expected error, got %q", tt.in, o)
				}
				return
			}
			if err != nil {
				t.Fatalf("Set(%q) unexpected error: %v", tt.in, err)
			}
			if o != tt.want {
				t.Errorf("Set(%q) = %q, want %q", tt.in, o, tt.want)
			}
		})
	}
}
