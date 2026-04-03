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

    os.Setenv("SERVER_HOST", "example.com")
    defer os.Unsetenv("SERVER_HOST")

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
