package config

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Validation struct {
		UsernamePattern string `yaml:"username" envconfig:"VALIDATION_USERNAME_PATTERN"`
		PasswordPattern string `yaml:"password" envconfig:"VALIDATION_PASSWORD_PATTERN"`
	} `yaml:"validation"`
}
