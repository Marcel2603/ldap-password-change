package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ldapService "github.com/Marcel2603/ldap-password-change/internal/service/ldap"

	"github.com/go-ldap/ldap/v3"
)

type mockServiceOK struct{}

func (*mockServiceOK) SearchUser(_ string) (*ldap.Entry, error) {
	return &ldap.Entry{DN: "cn=test,dc=example,dc=com"}, nil
}
func (*mockServiceOK) ChangePassword(_, _, _, _ string) error { return nil }
func (*mockServiceOK) Ping() error                            { return nil }

type mockServiceDown struct{}

func (*mockServiceDown) SearchUser(_ string) (*ldap.Entry, error) {
	return &ldap.Entry{DN: "cn=test,dc=example,dc=com"}, nil
}
func (*mockServiceDown) ChangePassword(_, _, _, _ string) error { return nil }
func (*mockServiceDown) Ping() error                            { return errors.New("connection refused") }

func TestLivenessHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health/live", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(LivenessHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
}

func TestReadinessHandler(t *testing.T) {
	tests := []struct {
		name           string
		svc            ldapService.Service
		expectedStatus int
		expectedBody   string
	}{
		{"ldap reachable", &mockServiceOK{}, http.StatusOK, "ok"},
		{"ldap unreachable", &mockServiceDown{}, http.StatusServiceUnavailable, "unavailable"},
	}
	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/health/ready", nil)
			rr := httptest.NewRecorder()
			ReadinessHandler(tc.svc).ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected %d, got %d", tc.expectedStatus, rr.Code)
			}
			var resp Response
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if resp.Status != tc.expectedBody {
				t.Errorf("expected status %q, got %q", tc.expectedBody, resp.Status)
			}
		})
	}
}
