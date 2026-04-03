package ldap_test

import (
	"crypto/tls"
	"errors"
	"io"
	"ldap-password-change/cmd/config"
	"ldap-password-change/internal/service/ldap"
	"log/slog"
	"testing"

	ldapext "github.com/go-ldap/ldap/v3"
)

type mockConn struct{}

func (*mockConn) Bind(_, _ string) error { return nil }
func (*mockConn) Close() error           { return nil }
func (*mockConn) Search(_ *ldapext.SearchRequest) (*ldapext.SearchResult, error) {
	return &ldapext.SearchResult{
		Entries: []*ldapext.Entry{
			{DN: "cn=tester,ou=users,dc=example,dc=org"},
		},
	}, nil
}
func (*mockConn) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return &ldapext.PasswordModifyResult{}, nil
}

type mockConnNoUser struct{ mockConn }

func (*mockConnNoUser) Search(_ *ldapext.SearchRequest) (*ldapext.SearchResult, error) {
	return &ldapext.SearchResult{Entries: []*ldapext.Entry{}}, nil
}

type mockConnErrorOnPwModify struct{ mockConn }

func (*mockConnErrorOnPwModify) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return nil, errors.New("test error on PasswordModify")
}

type mockConnErrorOnBind struct{ mockConn }

func (*mockConnErrorOnBind) Bind(_, _ string) error { return errors.New("test error on Bind") }

type mockLdapWrapperDefault struct{}

func (*mockLdapWrapperDefault) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConn{}, nil
}
func (*mockLdapWrapperDefault) DialWithTLSConfig(_ *tls.Config) ldapext.DialOpt {
	return func(_ *ldapext.DialContext) {}
}

type mockLdapWrapperNoUser struct{ mockLdapWrapperDefault }

func (*mockLdapWrapperNoUser) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnNoUser{}, nil
}

type mockLdapWrapperErrorOnPwModify struct{ mockLdapWrapperDefault }

func (*mockLdapWrapperErrorOnPwModify) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnErrorOnPwModify{}, nil
}

type mockLdapWrapperErrorOnCreate struct{ mockLdapWrapperDefault }

func (*mockLdapWrapperErrorOnCreate) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return nil, errors.New("test error on DialURL")
}

type mockLdapWrapperErrorOnBind struct{ mockLdapWrapperDefault }

func (*mockLdapWrapperErrorOnBind) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnErrorOnBind{}, nil
}

var defaultConfig = &config.LdapConfig{
	BaseDn:       "ou=users,dc=example,dc=org",
	UserDn:       "cn=admin,dc=example,dc=org",
	Password:     "123456",
	Host:         "unit.test:1389",
	IgnoreTLS:    true,
	SearchFilter: "(objectClass=*)",
}

func Test_serviceImpl_ChangePassword(t *testing.T) {
	tests := []struct {
		name            string
		config          config.LdapConfig
		ldapWrapperMock ldap.Wrapper
		username        string
		currentPassword string
		newPassword     string
		wantErr         bool
	}{
		{
			name:            "success",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperDefault{},
			username:        "tester",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         false,
		},
		{
			name:            "user not found",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperNoUser{},
			username:        "ghost",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         true,
		},
		{
			name:            "dial fails",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperErrorOnCreate{},
			username:        "tester",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         true,
		},
		{
			name:            "bind fails (wrong credentials)",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperErrorOnBind{},
			username:        "tester",
			currentPassword: "wrong",
			newPassword:     "Test1234",
			wantErr:         true,
		},
		{
			name:            "password modify fails",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperErrorOnPwModify{},
			username:        "tester",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := ldap.CreateService(tt.config, tt.ldapWrapperMock, mockLogger)
			err := svc.ChangePassword(tt.username, tt.currentPassword, tt.newPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_serviceImpl_Ping(t *testing.T) {
	mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ping succeeds", func(t *testing.T) {
		svc := ldap.CreateService(*defaultConfig, &mockLdapWrapperDefault{}, mockLogger)
		if err := svc.Ping(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("ping fails when dial fails", func(t *testing.T) {
		svc := ldap.CreateService(*defaultConfig, &mockLdapWrapperErrorOnCreate{}, mockLogger)
		if err := svc.Ping(); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ping fails when bind fails", func(t *testing.T) {
		svc := ldap.CreateService(*defaultConfig, &mockLdapWrapperErrorOnBind{}, mockLogger)
		if err := svc.Ping(); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
