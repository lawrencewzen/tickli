package utils

import (
	"testing"
	"time"

	"github.com/sho0pi/tickli/internal/config"
)

func TestShouldRefresh(t *testing.T) {
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

	base := func() *config.TokenData {
		return &config.TokenData{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ClientID:     "cid",
			ClientSecret: "secret",
			ExpiresAt:    now.Add(time.Hour),
		}
	}

	tests := []struct {
		name string
		mut  func(td *config.TokenData)
		want bool
	}{
		{
			name: "valid token far from expiry",
			mut:  func(td *config.TokenData) { td.ExpiresAt = now.Add(time.Hour) },
			want: false,
		},
		{
			name: "expires in 10 minutes",
			mut:  func(td *config.TokenData) { td.ExpiresAt = now.Add(10 * time.Minute) },
			want: false,
		},
		{
			name: "expires within threshold (4 min)",
			mut:  func(td *config.TokenData) { td.ExpiresAt = now.Add(4 * time.Minute) },
			want: true,
		},
		{
			name: "already expired",
			mut:  func(td *config.TokenData) { td.ExpiresAt = now.Add(-time.Minute) },
			want: true,
		},
		{
			name: "no refresh token",
			mut:  func(td *config.TokenData) { td.RefreshToken = "" },
			want: false,
		},
		{
			name: "no client id",
			mut:  func(td *config.TokenData) { td.ClientID = "" },
			want: false,
		},
		{
			name: "zero expiry (legacy token)",
			mut:  func(td *config.TokenData) { td.ExpiresAt = time.Time{} },
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := base()
			tt.mut(td)
			if got := shouldRefresh(td, now); got != tt.want {
				t.Errorf("shouldRefresh() = %v, want %v", got, tt.want)
			}
		})
	}
}
