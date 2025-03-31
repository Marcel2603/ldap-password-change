package validation_test

import (
	"ldap-password-change/cmd/config"
	"ldap-password-change/internal/validation"
	"testing"
)

const defaultPattern = "^[a-zA-Z0-9]*$"
const pcrePattern = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$"

type validationArgs struct {
	input   string
	pattern string
}

type validationTestCase struct {
	name string
	args validationArgs
	want bool
}

var validationTests = []validationTestCase{
	{
		name: "valid only letters",
		args: validationArgs{input: "test", pattern: defaultPattern},
		want: true,
	},
	{
		name: "invalid special characters",
		args: validationArgs{input: "ínvälid", pattern: defaultPattern},
		want: false,
	},
	{
		name: "valid only numbers",
		args: validationArgs{input: "1234", pattern: defaultPattern},
		want: true,
	},
	{
		name: "valid special characters",
		args: validationArgs{input: "test@12.3", pattern: "^[a-zA-Z0-9@\\.]*$"},
		want: true,
	},
	{
		name: "valid pcre pattern",
		args: validationArgs{input: "Test1234!", pattern: pcrePattern},
		want: true,
	},
	{
		name: "invalid pcre pattern",
		args: validationArgs{input: "test12345", pattern: pcrePattern},
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
			v, _ := validation.CreateValidator(validationConfig)
			if got := v.ValidateUsername(tt.args.input); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for %v with pattern %s", got, tt.want, tt.args.input, tt.args.pattern)
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
			v, _ := validation.CreateValidator(validationConfig)
			if got := v.ValidatePassword(tt.args.input); got != tt.want {
				t.Errorf("ValidateUsername() = %v, want %v for %v with pattern %s", got, tt.want, tt.args.input, tt.args.pattern)
			}
		})
	}
}

type configPatterns struct {
	usernamePattern string
	passwordPattern string
}

type creationTestCase struct {
	name    string
	args    configPatterns
	wantErr bool
}

var creationTests = []creationTestCase{
	{
		name:    "valid patterns",
		args:    configPatterns{usernamePattern: defaultPattern, passwordPattern: pcrePattern},
		wantErr: false,
	},
	{
		name:    "invalid username pattern",
		args:    configPatterns{usernamePattern: "?[", passwordPattern: defaultPattern},
		wantErr: true,
	},
	{
		name:    "invalid password pattern",
		args:    configPatterns{usernamePattern: defaultPattern, passwordPattern: "\\"},
		wantErr: true,
	},
}

func TestCreateValidator(t *testing.T) {
	for _, tt := range creationTests {
		t.Run(tt.name, func(t *testing.T) {
			validationConfig := config.ValidationConfig{
				UsernamePattern: tt.args.usernamePattern,
				PasswordPattern: tt.args.passwordPattern,
			}
			_, err := validation.CreateValidator(validationConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
