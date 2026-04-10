package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(cwd) }()

	_ = os.WriteFile("app.default.yml", []byte(`
server:
  host: localhost
  port: 8080
`), 0644)

	_ = os.WriteFile("app.yml", []byte(`
server:
  port: "8081"
`), 0644)

	_ = os.Setenv("SERVER_HOST", "example.com")
	defer func() {
		_ = os.Unsetenv("SERVER_HOST")
	}()

	loadConfig()

	// check if loaded
	if Configuration.Server.Port != "8081" {
		t.Errorf("Expected port 8081, got %v", Configuration.Server.Port)
	}
	if Configuration.Server.Host != "example.com" {
		t.Errorf("Expected host example.com, got %v", Configuration.Server.Host)
	}
}

func TestLoadConfigNoAppFile(t *testing.T) {
	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(cwd) }()

	_ = os.WriteFile("app.default.yml", []byte(`
server:
  host: localhost
  port: 8080
`), 0644)

	loadConfig()

	if Configuration.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %v", Configuration.Server.Port)
	}
}

func TestFormatUIConfig(t *testing.T) {
	c := Config{
		UI: UIConfig{
			BackgroundImage: "bg.jpg",
			CustomCSS:       "custom/mycss.css",
			Favicon:         "/custom/favicon.ico",
			Icon:            "https://example.com/logo.png",
		},
	}

	formatUIConfig(&c)

	if c.UI.BackgroundImage != "/custom/bg.jpg" {
		t.Errorf("Expected /custom/bg.jpg, got %v", c.UI.BackgroundImage)
	}
	if c.UI.CustomCSS != "/custom/mycss.css" {
		t.Errorf("Expected /custom/mycss.css, got %v", c.UI.CustomCSS)
	}
	if c.UI.Favicon != "/custom/favicon.ico" {
		t.Errorf("Expected /custom/favicon.ico, got %v", c.UI.Favicon)
	}
	if c.UI.Icon != "https://example.com/logo.png" {
		t.Errorf("Expected https://example.com/logo.png, got %v", c.UI.Icon)
	}

	cEmpty := Config{}
	formatUIConfig(&cEmpty)
	if cEmpty.UI.BackgroundImage != "" {
		t.Errorf("Expected empty string, got %v", cEmpty.UI.BackgroundImage)
	}
}
