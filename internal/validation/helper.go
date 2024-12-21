package validation

import (
	"github.com/dlclark/regexp2"
	"ldap-password-change/cmd/config"
)

var (
	usernameValidator = regexp2.MustCompile(config.Configuration.Validation.UsernamePattern, regexp2.None)
	passwordValidator = regexp2.MustCompile(config.Configuration.Validation.PasswordPattern, regexp2.None)
)

func ValidateUsername(username string) bool {
	isValid, _ := usernameValidator.MatchString(username)
	return isValid
}

func ValidatePassword(password string) bool {
	isValid, _ := passwordValidator.MatchString(password)
	return isValid
}
