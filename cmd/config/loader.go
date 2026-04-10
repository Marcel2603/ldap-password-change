package config

import (
	_ "embed"
	"log"
	"os"
	"strings"

	"dario.cat/mergo"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed app.default.yml
	yamlDefaultData []byte
	Configuration   Config
)

func init() {
	loadConfig()
}

func loadConfig() {
	var defaultData Config
	loadConfigFromYaml(yamlDefaultData, &defaultData)

	var appData Config
	appConfig, err := os.ReadFile("app.yml")
	if err == nil {
		log.Println("Loading app config from app.yml")
		loadConfigFromYaml(appConfig, &appData)
		_ = mergo.Merge(&defaultData, appData, mergo.WithOverride)
	}
	var envData Config
	loadConfigFromEnv(&envData)

	_ = mergo.Merge(&defaultData, envData, mergo.WithOverride)
	formatUIConfig(&defaultData)
	Configuration = defaultData
}

func formatUIConfig(c *Config) {
	prefix := func(path string) string {
		if path == "" || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "/") {
			return path
		}
		if strings.HasPrefix(path, "custom/") {
			return "/" + path
		}
		return "/custom/" + path
	}

	c.UI.BackgroundImage = prefix(c.UI.BackgroundImage)
	c.UI.CustomCSS = prefix(c.UI.CustomCSS)
	c.UI.Favicon = prefix(c.UI.Favicon)
	c.UI.Icon = prefix(c.UI.Icon)
}

func loadConfigFromEnv(mapData *Config) {
	err := envconfig.Process("", mapData)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func loadConfigFromYaml(data []byte, mapData *Config) {
	err := yaml.Unmarshal(data, &mapData)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}
