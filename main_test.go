package main

import (
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Marcel2603/ldap-password-change/cmd/config"
	staticfiles "github.com/Marcel2603/ldap-password-change/internal/handler/static-files"
	"github.com/Marcel2603/ldap-password-change/internal/service/ldap"

	goldap "github.com/go-ldap/ldap/v3"
)

type mockConn struct{}

func (*mockConn) Bind(_, _ string) error { return nil }
func (*mockConn) Close() error           { return nil }
func (*mockConn) Search(_ *goldap.SearchRequest) (*goldap.SearchResult, error) {
	return &goldap.SearchResult{
		Entries: []*goldap.Entry{
			{DN: "cn=test,ou=users,dc=example,dc=org"},
		},
	}, nil
}
func (*mockConn) PasswordModify(_ *goldap.PasswordModifyRequest) (*goldap.PasswordModifyResult, error) {
	return nil, nil
}

type mockWrapper struct{}

func (*mockWrapper) DialURL(_ string, _ ...goldap.DialOpt) (ldap.Conn, error) {
	return &mockConn{}, nil
}
func (*mockWrapper) DialWithTLSConfig(_ *tls.Config) goldap.DialOpt {
	return nil
}

func TestServerStarts(t *testing.T) {
	c := config.Config{}
	c.Server.Host = "localhost"
	c.Validation = config.ValidationConfig{
		UsernamePattern: "^.*$",
		PasswordPattern: "^.*$",
	}

	mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	staticfiles.NewHandler(staticFiles)

	app, err := setupApp(c, &mockWrapper{}, mockLogger)
	if err != nil {
		t.Fatalf("Failed to setup app: %v", err)
	}

	ts := httptest.NewServer(app)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}
