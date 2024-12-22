package validation_test

import (
	"ldap-password-change/cmd/config"
	"ldap-password-change/internal/validation"
	"testing"
)

type validationArgs struct {
	value   string
	pattern string
}

const defaultPattern = "^[a-zA-Z0-9]*$"

var validationTests = []struct {
	name string
	args validationArgs
	want bool
}{
	{
		name: "valid",
		args: validationArgs{value: "test", pattern: defaultPattern},
		want: true,
	},
	{
		name: "invalid",
		args: validationArgs{value: "1234", pattern: defaultPattern},
		want: true,
	},
	{
		name: "special characters",
		args: validationArgs{value: "test@12.3", pattern: "^[a-zA-Z0-9@\\.]*$"},
		want: true,
	},
	{
		name: "pcre pattern",
		args: validationArgs{value: "Test1234!", pattern: "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$"},
		want: true,
	},
}

func TestValidateUsername(t *testing.T) {
	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			config.Configuration.Validation.UsernamePattern = tt.args.pattern
			if got := validation.ValidateUsername(tt.args.value); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for pattern %s", got, tt.want, tt.args.pattern)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			config.Configuration.Validation.PasswordPattern = tt.args.pattern
			if got := validation.ValidatePassword(tt.args.value); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for pattern %s", got, tt.want, tt.args.pattern)
			}
		})
	}
}
