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

type validationTestCase struct {
	name string
	args validationArgs
	want bool
}

const defaultPattern = "^[a-zA-Z0-9]*$"
const pcrePattern = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$"

var validationTests = []validationTestCase{
	{
		name: "valid only letters",
		args: validationArgs{value: "test", pattern: defaultPattern},
		want: true,
	},
	{
		name: "invalid special characters",
		args: validationArgs{value: "ínvälid", pattern: defaultPattern},
		want: false,
	},
	{
		name: "valid only numbers",
		args: validationArgs{value: "1234", pattern: defaultPattern},
		want: true,
	},
	{
		name: "valid special characters",
		args: validationArgs{value: "test@12.3", pattern: "^[a-zA-Z0-9@\\.]*$"},
		want: true,
	},
	{
		name: "valid pcre pattern",
		args: validationArgs{value: "Test1234!", pattern: pcrePattern},
		want: true,
	},
	{
		name: "invalid pcre pattern",
		args: validationArgs{value: "test12345", pattern: pcrePattern},
		want: false,
	},
}

func TestValidateUsername(t *testing.T) {
	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			validationConfig := config.ValidationConfig{
				UsernamePattern: tt.args.pattern,
				PasswordPattern: "^.*$",
			}
			v := validation.CreateValidator(validationConfig)
			if got := v.ValidateUsername(tt.args.value); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for %v with pattern %s", got, tt.want, tt.args.value, tt.args.pattern)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			validationConfig := config.ValidationConfig{
				UsernamePattern: "^.*$",
				PasswordPattern: tt.args.pattern,
			}
			v := validation.CreateValidator(validationConfig)
			if got := v.ValidatePassword(tt.args.value); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for %v with pattern %s", got, tt.want, tt.args.value, tt.args.pattern)
			}
		})
	}
}
