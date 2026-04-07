package ldap_test

import (
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/Marcel2603/ldap-password-change/cmd/config"
	"github.com/Marcel2603/ldap-password-change/internal/service/ldap"

	ldapext "github.com/go-ldap/ldap/v3"
)

type mockConn struct {
	BindFunc           func(username, password string) error
	SearchFunc         func(searchRequest *ldapext.SearchRequest) (*ldapext.SearchResult, error)
	PasswordModifyFunc func(passwordModifyRequest *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error)
	CloseFunc          func() error
}

func (m *mockConn) Bind(u, p string) error {
	if m.BindFunc != nil {
		return m.BindFunc(u, p)
	}
	return nil
}

func (m *mockConn) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *mockConn) Search(req *ldapext.SearchRequest) (*ldapext.SearchResult, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(req)
	}
	return &ldapext.SearchResult{
		Entries: []*ldapext.Entry{{DN: "cn=tester,ou=users,dc=example,dc=org"}},
	}, nil
}

func (m *mockConn) PasswordModify(req *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	if m.PasswordModifyFunc != nil {
		return m.PasswordModifyFunc(req)
	}
	return &ldapext.PasswordModifyResult{}, nil
}

type mockWrapper struct {
	DialURLFunc func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error)
}

func (m *mockWrapper) DialURL(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
	if m.DialURLFunc != nil {
		return m.DialURLFunc(addr, opts...)
	}
	return &mockConn{}, nil
}

func (m *mockWrapper) DialWithTLSConfig(_ *tls.Config) ldapext.DialOpt {
	return func(_ *ldapext.DialContext) {}
}

var defaultConfig = &config.LdapConfig{
	BaseDn:       "ou=users,dc=example,dc=org",
	UserDn:       "cn=admin,dc=example,dc=org",
	Password:     "123456",
	Host:         "unit.test:1389",
	IgnoreTLS:    true,
	SearchFilter: "(objectClass=*)",
}

func Test_serviceImpl_SearchUser(t *testing.T) {
	tests := []struct {
		name            string
		config          config.LdapConfig
		ldapWrapperMock ldap.Wrapper
		username        string
		wantDN          string
		wantErr         bool
	}{
		{
			name:            "success",
			config:          *defaultConfig,
			ldapWrapperMock: &mockWrapper{},
			username:        "tester",
			wantDN:          "cn=tester,ou=users,dc=example,dc=org",
			wantErr:         false,
		},
		{
			name:   "user not found",
			config: *defaultConfig,
			ldapWrapperMock: &mockWrapper{
				DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
					return &mockConn{
						SearchFunc: func(req *ldapext.SearchRequest) (*ldapext.SearchResult, error) {
							return &ldapext.SearchResult{Entries: []*ldapext.Entry{}}, nil
						},
					}, nil
				},
			},
			username: "ghost",
			wantErr:  true,
		},
		{
			name:   "dial fails",
			config: *defaultConfig,
			ldapWrapperMock: &mockWrapper{
				DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
					return nil, errors.New("test error on DialURL")
				},
			},
			username: "tester",
			wantErr:  true,
		},
		{
			name:   "search returns error",
			config: *defaultConfig,
			ldapWrapperMock: &mockWrapper{
				DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
					return &mockConn{
						SearchFunc: func(req *ldapext.SearchRequest) (*ldapext.SearchResult, error) {
							return nil, errors.New("search failed")
						},
					}, nil
				},
			},
			username: "tester",
			wantErr:  true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := ldap.CreateService(tt.config, tt.ldapWrapperMock, mockLogger)

			entry, err := svc.SearchUser(tt.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && entry.DN != tt.wantDN {
				t.Errorf("SearchUser() got DN = %v, want %v", entry.DN, tt.wantDN)
			}
		})
	}
}

func Test_serviceImpl_ChangePassword(t *testing.T) {
	tests := []struct {
		name            string
		config          config.LdapConfig
		ldapWrapperMock ldap.Wrapper
		userDN          string
		username        string
		currentPassword string
		newPassword     string
		wantErr         bool
	}{
		{
			name:            "success",
			config:          *defaultConfig,
			ldapWrapperMock: &mockWrapper{},
			userDN:          "cn=tester,ou=users,dc=example,dc=org",
			username:        "tester",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         false,
		},
		{
			name:   "dial fails",
			config: *defaultConfig,
			ldapWrapperMock: &mockWrapper{
				DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
					return nil, errors.New("test error on DialURL")
				},
			},
			userDN:          "cn=tester,ou=users,dc=example,dc=org",
			username:        "tester",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         true,
		},
		{
			name:   "bind fails (wrong credentials)",
			config: *defaultConfig,
			ldapWrapperMock: &mockWrapper{
				DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
					return &mockConn{
						BindFunc: func(u, p string) error { return errors.New("test error on Bind") },
					}, nil
				},
			},
			userDN:          "cn=tester,ou=users,dc=example,dc=org",
			username:        "tester",
			currentPassword: "wrong",
			newPassword:     "Test1234",
			wantErr:         true,
		},
		{
			name:   "password modify fails",
			config: *defaultConfig,
			ldapWrapperMock: &mockWrapper{
				DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
					return &mockConn{
						PasswordModifyFunc: func(req *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
							return nil, errors.New("test error on PasswordModify")
						},
					}, nil
				},
			},
			userDN:          "cn=tester,ou=users,dc=example,dc=org",
			username:        "tester",
			currentPassword: "123456",
			newPassword:     "Test1234",
			wantErr:         true,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := ldap.CreateService(tt.config, tt.ldapWrapperMock, mockLogger)

			err := svc.ChangePassword(tt.userDN, tt.username, tt.currentPassword, tt.newPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_serviceImpl_Ping(t *testing.T) {
	mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ping succeeds", func(t *testing.T) {
		svc := ldap.CreateService(*defaultConfig, &mockWrapper{}, mockLogger)
		if err := svc.Ping(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("ping fails when dial fails", func(t *testing.T) {
		mockWrapperErr := &mockWrapper{
			DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
				return nil, errors.New("dial failed")
			},
		}
		svc := ldap.CreateService(*defaultConfig, mockWrapperErr, mockLogger)
		if err := svc.Ping(); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ping fails when bind fails", func(t *testing.T) {
		mockWrapperErr := &mockWrapper{
			DialURLFunc: func(addr string, opts ...ldapext.DialOpt) (ldap.Conn, error) {
				return &mockConn{
					BindFunc: func(u, p string) error { return errors.New("bind failed") },
				}, nil
			},
		}
		svc := ldap.CreateService(*defaultConfig, mockWrapperErr, mockLogger)
		if err := svc.Ping(); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
