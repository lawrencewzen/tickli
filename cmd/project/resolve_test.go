package project

import (
	"strings"
	"testing"

	"github.com/sho0pi/tickli/internal/types"
)

func projects() []types.Project {
	return []types.Project{
		{ID: "id-work", Name: "Work Tasks"},
		{ID: "id-home", Name: "Home"},
		{ID: "id-side", Name: "Side Work"},
	}
}

func TestResolveProject(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantID  string
		wantErr bool
		// errContains is checked (case-insensitively) against the error message
		// when wantErr is true.
		errContains string
	}{
		{
			name:   "exact id match",
			query:  "id-home",
			wantID: "id-home",
		},
		{
			name:   "exact id wins over name substring of others",
			query:  "id-work",
			wantID: "id-work",
		},
		{
			name:   "unique name substring, case-insensitive",
			query:  "HOME",
			wantID: "id-home",
		},
		{
			name:        "no match",
			query:       "nonexistent",
			wantErr:     true,
			errContains: "no project matches",
		},
		{
			name:        "ambiguous name substring",
			query:       "work",
			wantErr:     true,
			errContains: "multiple projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveProject(projects(), tt.query)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got project %+v", got)
				}
				if tt.errContains != "" && !strings.Contains(strings.ToLower(err.Error()), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.wantID {
				t.Errorf("resolveProject() id = %q, want %q", got.ID, tt.wantID)
			}
		})
	}
}

// The multi-match error should list every candidate so the user can disambiguate.
func TestResolveProject_AmbiguousListsCandidates(t *testing.T) {
	_, err := resolveProject(projects(), "work")
	if err == nil {
		t.Fatal("expected ambiguity error")
	}
	for _, want := range []string{"Work Tasks", "Side Work", "id-work", "id-side"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error message missing candidate %q:\n%s", want, err.Error())
		}
	}
}
