package ldap_test

import (
	"crypto/tls"
	"errors"
	ldapext "github.com/go-ldap/ldap/v3"
	"ldap-password-change/cmd/config"
	"ldap-password-change/internal/service/ldap"
	"testing"
)

type changePasswordArgs struct {
	username        string
	currentPassword string
	newPassword     string
}

type changePasswordTestCase struct {
	name            string
	config          config.LdapConfig
	ldapWrapperMock ldap.LdapWrapper
	args            changePasswordArgs
	wantErr         bool
}

type mockConn struct {
}

func (l *mockConn) Bind(_, _ string) error {
	return nil
}
func (l *mockConn) Close() error {
	return nil
}
func (l *mockConn) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return &ldapext.PasswordModifyResult{}, nil
}

type mockConnError struct {
	mockConn
}

func (l *mockConnError) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return nil, errors.New("test error")
}

type mockLdapWrapperDefault struct {
}

func (w *mockLdapWrapperDefault) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConn{}, nil
}
func (w *mockLdapWrapperDefault) DialWithTLSConfig(_ *tls.Config) ldapext.DialOpt {
	return func(dc *ldapext.DialContext) {}
}

type mockLdapWrapperError struct {
	mockLdapWrapperDefault
}

func (w *mockLdapWrapperError) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnError{}, nil
}

var (
	defaultConfig = &config.LdapConfig{
		BaseDn:   "ou=users,dc=example,dc=org",
		UserDn:   "cn=admin,dc=example,dc=org",
		Password: "123456",
		Domain:   "ldap://unit.test:1389",
	}
)

func Test_serviceImpl_ChangePassword(t *testing.T) {
	tests := []changePasswordTestCase{
		{
			name:            "change password should succeed",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperDefault{},
			args: changePasswordArgs{
				username:        "tester",
				currentPassword: "123456",
				newPassword:     "Test1234",
			},
			wantErr: false,
		},
		{
			name:            "change password should fail",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperError{},
			args: changePasswordArgs{
				username:        "tester",
				currentPassword: "123456",
				newPassword:     "Test1234",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ldap.CreateService(tt.config, tt.ldapWrapperMock)
			if err := s.ChangePassword(tt.args.username, tt.args.currentPassword, tt.args.newPassword); (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
