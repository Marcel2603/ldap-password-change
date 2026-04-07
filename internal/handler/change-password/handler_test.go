package changepassword

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

type mockValidator struct {
	usernameValid bool
	passwordValid bool
}

func (m *mockValidator) ValidateUsername(username string) bool { return m.usernameValid }
func (m *mockValidator) ValidatePassword(password string) bool { return m.passwordValid }

type mockLdapService struct {
	returnsError error
	username     string
	oldPwd       string
	newPwd       string
}

func (m *mockLdapService) SearchUser(username string) (*ldap.Entry, error) {
	return &ldap.Entry{DN: "cn=testuser,dc=example,dc=com"}, nil
}

func (m *mockLdapService) ChangePassword(userDN string, username string, currentPassword string, newPassword string) error {
	m.username = username
	m.oldPwd = currentPassword
	m.newPwd = newPassword
	return m.returnsError
}

func (*mockLdapService) Ping() error { return nil }

func TestHandler(t *testing.T) {
	tests := []struct {
		name               string
		usernameValid      bool
		passwordValid      bool
		passwordsMatch     bool
		ldapError          error
		expectedStatusCode int
	}{
		{
			name:               "successful password change",
			usernameValid:      true,
			passwordValid:      true,
			passwordsMatch:     true,
			ldapError:          nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "invalid username",
			usernameValid:      false,
			passwordValid:      true,
			passwordsMatch:     true,
			ldapError:          nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid password",
			usernameValid:      true,
			passwordValid:      false,
			passwordsMatch:     true,
			ldapError:          nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "passwords do not match",
			usernameValid:      true,
			passwordValid:      true,
			passwordsMatch:     false,
			ldapError:          nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "ldap change fails",
			usernameValid:      true,
			passwordValid:      true,
			passwordsMatch:     true,
			ldapError:          errors.New("ldap error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ldapSvc := &mockLdapService{returnsError: tc.ldapError}
			validator := &mockValidator{usernameValid: tc.usernameValid, passwordValid: tc.passwordValid}

			mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := Handler(ldapSvc, validator, mockLogger)

			form := url.Values{}
			form.Add("username", "testuser")
			form.Add("current-password", "oldpass")
			form.Add("new-password", "newpass")
			if tc.passwordsMatch {
				form.Add("confirm-password", "newpass")
			} else {
				form.Add("confirm-password", "wrongpass")
			}

			req, err := http.NewRequest("POST", "/change-password", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tc.expectedStatusCode, rr.Code)
			}
		})
	}
}
