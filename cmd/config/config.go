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
	Host         string `yaml:"host"`
	UserDn       string `yaml:"userDn"`
	Password     string `yaml:"password"`
	BaseDn       string `yaml:"baseDn"`
	SearchFilter string `yaml:"searchFilter"`
	IgnoreTLS    bool   `yaml:"ignoreTLS"`
	TlsCert      string `yaml:"tlsCert"`
}

type ValidationConfig struct {
	UsernamePattern string `yaml:"username"`
	PasswordPattern string `yaml:"password"`
}
