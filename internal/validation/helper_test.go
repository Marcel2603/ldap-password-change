package validation_test

import (
	"github.com/dlclark/regexp2"
	"ldap-password-change/internal/validation"
	"testing"
)

type validationArgs struct {
	value   string
	pattern string
}

const defaultPattern = "^[a-zA-Z0-9]*$"
const pcrePattern = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$"

var validationTests = []struct {
	name string
	args validationArgs
	want bool
}{
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
			validation.UsernameValidator = regexp2.MustCompile(tt.args.pattern, regexp2.None)
			if got := validation.ValidateUsername(tt.args.value); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for %v with pattern %s", got, tt.want, tt.args.value, tt.args.pattern)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			validation.PasswordValidator = regexp2.MustCompile(tt.args.pattern, regexp2.None)
			if got := validation.ValidatePassword(tt.args.value); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for %v with pattern %s", got, tt.want, tt.args.value, tt.args.pattern)
			}
		})
	}
}
