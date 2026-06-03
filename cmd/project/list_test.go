package project

import (
	"testing"

	"github.com/sho0pi/tickli/internal/types"
)

func TestFilterProjectByName(t *testing.T) {
	all := []types.Project{
		{ID: "1", Name: "Work Tasks"},
		{ID: "2", Name: "Home"},
		{ID: "3", Name: "Side Work"},
	}

	tests := []struct {
		name    string
		query   string
		wantIDs []string
		wantErr bool
	}{
		{name: "single match", query: "home", wantIDs: []string{"2"}},
		{name: "multiple substring matches", query: "WORK", wantIDs: []string{"1", "3"}},
		{name: "no match returns error", query: "zzz", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterProjectByName(all, tt.query)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("got %d projects, want %d (%+v)", len(got), len(tt.wantIDs), got)
			}
			for i, id := range tt.wantIDs {
				if got[i].ID != id {
					t.Errorf("match[%d] id = %q, want %q", i, got[i].ID, id)
				}
			}
		})
	}
}
