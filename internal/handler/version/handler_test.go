package version_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Marcel2603/ldap-password-change/internal/handler/version"
)

func TestVersionHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/version", nil)
	rr := httptest.NewRecorder()
	versionStr := "1.0.0"
	commit := "1234ef"
	version.Handler(versionStr, commit).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp version.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Version != versionStr {
		t.Errorf("expected version '%s', got %q", versionStr, resp.Version)
	}
	if resp.Commit != commit {
		t.Errorf("expected commit '%s', got %q", commit, resp.Commit)
	}
}
