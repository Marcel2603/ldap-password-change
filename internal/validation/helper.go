package validation

import (
	"errors"
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

func CreateValidator(config config.ValidationConfig) (Validator, error) {
	usernameValidatorRegexp, err1 := regexp2.Compile(config.UsernamePattern, regexp2.None)
	passwortValidatorRegexp, err2 := regexp2.Compile(config.PasswordPattern, regexp2.None)
	if err1 != nil || err2 != nil {
		return nil, errors.Join(err1, err2)
	}

	return &validator{
		usernameValidator: usernameValidatorRegexp,
		passwordValidator: passwortValidatorRegexp,
	}, nil
}

func (v *validator) ValidateUsername(username string) bool {
	isValid, _ := v.usernameValidator.MatchString(username)
	return isValid
}

func (v *validator) ValidatePassword(password string) bool {
	isValid, _ := v.passwordValidator.MatchString(password)
	return isValid
}
