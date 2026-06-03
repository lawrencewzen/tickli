package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sho0pi/tickli/internal/types"
)

// newTestClient builds a Client pointed at a mock server. It bypasses NewClient
// so tests don't talk to the real TickTick base URL.
func newTestClient(baseURL string) *Client {
	return &Client{http: resty.New().SetBaseURL(baseURL)}
}

func TestGetAuthURL(t *testing.T) {
	got := GetAuthURL("my-client-id")
	for _, want := range []string{"client_id=my-client-id", "scope=" + scope, "response_type=code", redirectURL} {
		if !strings.Contains(got, want) {
			t.Errorf("auth URL %q missing %q", got, want)
		}
	}
}

func TestListProjects_AppendsInbox(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/project" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"p1","name":"Work"}]`))
	}))
	defer srv.Close()

	got, err := newTestClient(srv.URL).ListProjects()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 projects (api + inbox), got %d: %+v", len(got), got)
	}
	if got[len(got)-1].ID != types.InboxProject.ID {
		t.Errorf("expected last project to be inbox, got %q", got[len(got)-1].ID)
	}
}

func TestListProjects_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`boom`))
	}))
	defer srv.Close()

	if _, err := newTestClient(srv.URL).ListProjects(); err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestListTasks_UnwrapsTasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/project/pid/data"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tasks":[{"id":"t1","title":"first"},{"id":"t2","title":"second"}]}`))
	}))
	defer srv.Close()

	got, err := newTestClient(srv.URL).ListTasks("pid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].ID != "t1" || got[1].Title != "second" {
		t.Errorf("unexpected tasks: %+v", got)
	}
}

func TestGetTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/project/pid/task/tid"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"tid","title":"hello"}`))
	}))
	defer srv.Close()

	got, err := newTestClient(srv.URL).GetTask("pid", "tid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "tid" || got.Title != "hello" {
		t.Errorf("unexpected task: %+v", got)
	}
}

func TestCreateTask(t *testing.T) {
	t.Run("nil task is rejected without a request", func(t *testing.T) {
		called := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		}))
		defer srv.Close()

		if _, err := newTestClient(srv.URL).CreateTask(nil); err == nil {
			t.Fatal("expected error for nil task")
		}
		if called {
			t.Error("server should not be called for nil task")
		}
	})

	t.Run("returns server-assigned task", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != "/task" {
				t.Errorf("got %s %s, want POST /task", r.Method, r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"generated-id","title":"Buy milk"}`))
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).CreateTask(&types.Task{Title: "Buy milk"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != "generated-id" {
			t.Errorf("ID = %q, want generated-id", got.ID)
		}
	})
}

func TestCompleteTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/project/pid/task/tid/complete"; got != want {
				t.Errorf("path = %q, want %q", got, want)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		if err := newTestClient(srv.URL).CompleteTask("pid", "tid"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("server error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer srv.Close()

		if err := newTestClient(srv.URL).CompleteTask("pid", "tid"); err == nil {
			t.Fatal("expected error on 403 response")
		}
	})
}

func TestGetAccessToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		if got := r.FormValue("grant_type"); got != "authorization_code" {
			t.Errorf("grant_type = %q, want authorization_code", got)
		}
		if got := r.FormValue("code"); got != "the-code" {
			t.Errorf("code = %q, want the-code", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"at","refresh_token":"rt","expires_in":3600}`))
	}))
	defer srv.Close()
	withTokenURL(t, srv.URL)

	resp, err := GetAccessToken("cid", "secret", "the-code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "at" || resp.RefreshToken != "rt" || resp.ExpiresIn != 3600 {
		t.Errorf("unexpected token response: %+v", resp)
	}
}
