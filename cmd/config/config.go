package config

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Validation ValidationConfig `yaml:"validation"`
	Ldap       LdapConfig       `yaml:"ldap"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type LdapConfig struct {
	Domain       string `yaml:"domain"`
	UserDn       string `yaml:"userDn"`
	Password     string `yaml:"password"`
	BaseDn       string `yaml:"baseDn"`
	SearchFilter string `yaml:"searchFilter"`
}

type ValidationConfig struct {
	UsernamePattern string `yaml:"username"`
	PasswordPattern string `yaml:"password"`
}
