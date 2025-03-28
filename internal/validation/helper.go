package validation

import (
	"github.com/dlclark/regexp2"
	"ldap-password-change/cmd/config"
)

type Validator interface {
	ValidateUsername(username string) bool
	ValidatePassword(password string) bool
}

type validator struct {
	usernameValidator *regexp2.Regexp
	passwordValidator *regexp2.Regexp
}

func CreateValidator(config config.ValidationConfig) Validator {
	return &validator{
		usernameValidator: regexp2.MustCompile(config.UsernamePattern, regexp2.None),
		passwordValidator: regexp2.MustCompile(config.PasswordPattern, regexp2.None),
	}
}

func (v *validator) ValidateUsername(username string) bool {
	isValid, _ := v.usernameValidator.MatchString(username)
	return isValid
}

func (v *validator) ValidatePassword(password string) bool {
	isValid, _ := v.passwordValidator.MatchString(password)
	return isValid
}
