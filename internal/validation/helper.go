package validation

import (
	"github.com/dlclark/regexp2"
	"ldap-password-change/cmd/config"
)

type Validator struct {
	usernameValidator *regexp2.Regexp
	passwordValidator *regexp2.Regexp
}

func CreateValidator(config config.Config) Validator {
	return Validator{
		usernameValidator: regexp2.MustCompile(config.Validation.UsernamePattern, regexp2.None),
		passwordValidator: regexp2.MustCompile(config.Validation.PasswordPattern, regexp2.None),
	}
}

func (v *Validator) ValidateUsername(username string) bool {
	isValid, _ := v.usernameValidator.MatchString(username)
	return isValid
}

func (v *Validator) ValidatePassword(password string) bool {
	isValid, _ := v.passwordValidator.MatchString(password)
	return isValid
}
