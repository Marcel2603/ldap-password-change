package config

type Config struct {
  Server struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
  } `yaml:"server"`
  Validation struct {
    UsernamePattern string `yaml:"username"`
    PasswordPattern string `yaml:"password"`
  } `yaml:"validation"`
}
