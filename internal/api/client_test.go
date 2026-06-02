package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// withTokenURL points tokenURL at a mock server for the duration of a test.
func withTokenURL(t *testing.T, url string) {
	t.Helper()
	orig := tokenURL
	tokenURL = url
	t.Cleanup(func() { tokenURL = orig })
}

func TestRefreshAccessToken_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if got := r.FormValue("grant_type"); got != "refresh_token" {
			t.Errorf("grant_type = %q, want refresh_token", got)
		}
		if got := r.FormValue("refresh_token"); got != "old-refresh" {
			t.Errorf("refresh_token = %q, want old-refresh", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh","expires_in":3600,"token_type":"bearer"}`))
	}))
	defer srv.Close()
	withTokenURL(t, srv.URL)

	resp, err := RefreshAccessToken("cid", "secret", "old-refresh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want new-access", resp.AccessToken)
	}
	if resp.RefreshToken != "new-refresh" {
		t.Errorf("RefreshToken = %q, want new-refresh", resp.RefreshToken)
	}
	if resp.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %d, want 3600", resp.ExpiresIn)
	}
}

func TestRefreshAccessToken_Rejected(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer srv.Close()
	withTokenURL(t, srv.URL)

	resp, err := RefreshAccessToken("cid", "secret", "expired-refresh")
	if err == nil {
		t.Fatalf("expected error for rejected refresh token, got resp %+v", resp)
	}
	if resp != nil {
		t.Errorf("expected nil response on error, got %+v", resp)
	}
}
