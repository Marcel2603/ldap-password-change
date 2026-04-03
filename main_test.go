package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	goldap "github.com/go-ldap/ldap/v3"
	"ldap-password-change/cmd/config"
	"ldap-password-change/internal/service/ldap"
)

type mockConn struct{}

func (m *mockConn) Bind(username, password string) error { return nil }
func (m *mockConn) Close() error                         { return nil }
func (m *mockConn) PasswordModify(req *goldap.PasswordModifyRequest) (*goldap.PasswordModifyResult, error) {
	return nil, nil
}

type mockWrapper struct{}

func (m *mockWrapper) DialURL(addr string, opts ...goldap.DialOpt) (ldap.Conn, error) {
	return &mockConn{}, nil
}
func (m *mockWrapper) DialWithTLSConfig(tc *tls.Config) goldap.DialOpt {
	return nil
}

func TestServerStarts(t *testing.T) {
	c := config.Config{}
	c.Server.Host = "localhost"
	c.Validation = config.ValidationConfig{
		UsernamePattern: "^.*$",
		PasswordPattern: "^.*$",
	}

	app, err := setupApp(c, &mockWrapper{})
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
