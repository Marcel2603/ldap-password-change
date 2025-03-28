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

// TODO #3: this implementation of mocks is quite messy, surely there is a better way

type mockConnOk struct {
}

func (l *mockConnOk) Bind(_, _ string) error {
	return nil
}
func (l *mockConnOk) Close() error {
	return nil
}
func (l *mockConnOk) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return &ldapext.PasswordModifyResult{}, nil
}

type mockConnError struct {
}

func (l *mockConnError) Bind(_, _ string) error {
	return nil
}
func (l *mockConnError) Close() error {
	return nil
}
func (l *mockConnError) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return nil, errors.New("test error")
}

type ldapWrapperOk struct {
}

func (w *ldapWrapperOk) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnOk{}, nil
}
func (w *ldapWrapperOk) DialWithTLSConfig(_ *tls.Config) ldapext.DialOpt {
	return func(dc *ldapext.DialContext) {}
}

type ldapWrapperError struct {
}

func (w *ldapWrapperError) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnError{}, nil
}
func (w *ldapWrapperError) DialWithTLSConfig(_ *tls.Config) ldapext.DialOpt {
	return func(dc *ldapext.DialContext) {}
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
			ldapWrapperMock: &ldapWrapperOk{},
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
			ldapWrapperMock: &ldapWrapperError{},
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
