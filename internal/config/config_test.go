package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// withTempTokenPath points tokenPath at a throwaway file for the duration of a
// test and restores the original afterwards.
func withTempTokenPath(t *testing.T) {
	t.Helper()
	orig := tokenPath
	tokenPath = filepath.Join(t.TempDir(), "token")
	t.Cleanup(func() { tokenPath = orig })
}

func TestLoadTokenData_Missing(t *testing.T) {
	withTempTokenPath(t)

	td, err := LoadTokenData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if td != nil {
		t.Fatalf("expected nil token data for missing file, got %+v", td)
	}
}

func TestSaveLoadTokenData_RoundTrip(t *testing.T) {
	withTempTokenPath(t)

	want := &TokenData{
		AccessToken:  "access-123",
		RefreshToken: "refresh-456",
		ExpiresAt:    time.Now().Add(time.Hour).Truncate(time.Second),
		ClientID:     "cid",
		ClientSecret: "secret",
	}
	if err := SaveTokenData(want); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := LoadTokenData()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if got.AccessToken != want.AccessToken ||
		got.RefreshToken != want.RefreshToken ||
		got.ClientID != want.ClientID ||
		got.ClientSecret != want.ClientSecret {
		t.Errorf("round trip mismatch:\n got %+v\nwant %+v", got, want)
	}
	if !got.ExpiresAt.Equal(want.ExpiresAt) {
		t.Errorf("ExpiresAt mismatch: got %v want %v", got.ExpiresAt, want.ExpiresAt)
	}
}

// TestLoadTokenData_LegacyRawString covers backward compatibility with the old
// format, where the token file held a bare access token string rather than JSON.
func TestLoadTokenData_LegacyRawString(t *testing.T) {
	withTempTokenPath(t)

	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(tokenPath, []byte("legacy-raw-token"), 0600); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	td, err := LoadTokenData()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if td.AccessToken != "legacy-raw-token" {
		t.Errorf("expected legacy token to load as access token, got %q", td.AccessToken)
	}
	if td.RefreshToken != "" {
		t.Errorf("expected no refresh token for legacy format, got %q", td.RefreshToken)
	}
}
