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
	name              string
	config            config.LdapConfig
	ldapWrapperMock   ldap.LdapWrapper
	args              changePasswordArgs
	wantErrOnCreation bool
	wantErrOnAction   bool
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

type mockConnErrorOnPwModify struct {
	mockConn
}

func (l *mockConnErrorOnPwModify) PasswordModify(_ *ldapext.PasswordModifyRequest) (*ldapext.PasswordModifyResult, error) {
	return nil, errors.New("test error on PasswordModify")
}

type mockLdapWrapperDefault struct {
}

func (w *mockLdapWrapperDefault) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConn{}, nil
}
func (w *mockLdapWrapperDefault) DialWithTLSConfig(_ *tls.Config) ldapext.DialOpt {
	return func(dc *ldapext.DialContext) {}
}

type mockLdapWrapperErrorOnPwModify struct {
	mockLdapWrapperDefault
}

func (w *mockLdapWrapperErrorOnPwModify) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return &mockConnErrorOnPwModify{}, nil
}

type mockLdapWrapperErrorOnCreate struct {
	mockLdapWrapperDefault
}

func (w *mockLdapWrapperErrorOnCreate) DialURL(_ string, _ ...ldapext.DialOpt) (ldap.Conn, error) {
	return nil, errors.New("test error on DialURL")
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
			wantErrOnCreation: false,
			wantErrOnAction:   false,
		},
		{
			name:            "change password should fail when pw modify fails",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperErrorOnPwModify{},
			args: changePasswordArgs{
				username:        "tester",
				currentPassword: "123456",
				newPassword:     "Test1234",
			},
			wantErrOnCreation: false,
			wantErrOnAction:   true,
		},
		{
			name:            "change password should fail when client connection fails",
			config:          *defaultConfig,
			ldapWrapperMock: &mockLdapWrapperErrorOnCreate{},
			args: changePasswordArgs{
				username:        "tester",
				currentPassword: "123456",
				newPassword:     "Test1234",
			},
			wantErrOnCreation: true,
			wantErrOnAction:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, errCreation := ldap.CreateService(tt.config, tt.ldapWrapperMock)
			if errCreation != nil {
				if !tt.wantErrOnCreation {
					t.Errorf("CreateService() error = %v, wantErrOnCreation %v", errCreation, tt.wantErrOnCreation)
				}
			} else {
				if errAction := s.ChangePassword(tt.args.username, tt.args.currentPassword, tt.args.newPassword); (errAction != nil) != tt.wantErrOnAction {
					t.Errorf("ChangePassword() error = %v, wantErrOnAction %v", errAction, tt.wantErrOnAction)
				}
			}
		})
	}
}
