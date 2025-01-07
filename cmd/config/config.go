package config

type Config struct {
	Server     server     `yaml:"server"`
	Validation validation `yaml:"validation"`
	Ldap       ldap       `yaml:"ldap"`
}

type server struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type ldap struct {
	Domain     string `yaml:"domain"`
	UserDn     string `yaml:"userDn"`
	Password   string `yaml:"password"`
	BaseDn     string `yaml:"baseDn"`
	UserFilter string `yaml:"userFilter"`
}

type validation struct {
	UsernamePattern string `yaml:"username"`
	PasswordPattern string `yaml:"password"`
}
