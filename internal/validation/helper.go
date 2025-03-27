package validation

import (
  "github.com/dlclark/regexp2"
  "ldap-password-change/cmd/config"
)

var (
  UsernameValidator = regexp2.MustCompile(config.Configuration.Validation.UsernamePattern, regexp2.None)
  PasswordValidator = regexp2.MustCompile(config.Configuration.Validation.PasswordPattern, regexp2.None)
)

func ValidateUsername(username string) bool {
  isValid, _ := UsernameValidator.MatchString(username)
  return isValid
}

func ValidatePassword(password string) bool {
  isValid, _ := PasswordValidator.MatchString(password)
  return isValid
}
