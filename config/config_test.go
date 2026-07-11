package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig_LocalFile(t *testing.T) {
	configContent := `{"logLevel":"debug"}`

	if err := os.MkdirAll("config", 0o755); err != nil {
		t.Fatalf("could not create directory 'config': %v", err)
	}
	configPath := filepath.Join("config", "test_config.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("could not write config file: %v", err)
	}
	defer os.Remove(configPath) // nolint

	cfg, err := newConfig("test_config.json")
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Unexpected config: %+v", cfg)
	}
}

func TestNewConfig_HomeFallback(t *testing.T) {
	configContent := `{"logLevel":"info"}`

	tmpHome := t.TempDir()
	oldUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) { return tmpHome, nil }
	defer func() { userHomeDir = oldUserHomeDir }()

	configDir := filepath.Join(tmpHome, ".config", "nirimgr")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("could not create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "test_config.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("could not write config file: %v", err)
	}

	cfg, err := newConfig("test_config.json")
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Unexpected config: %+v", cfg)
	}
}

func TestConfigure_SetsGlobalConfig(t *testing.T) {
	configContent := `{"logLevel":"warn"}`

	if err := os.MkdirAll("config", 0o755); err != nil {
		t.Fatalf("could not create directory 'config': %v", err)
	}
	configPath := filepath.Join("config", "test_config.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("could not write config file: %v", err)
	}
	defer os.Remove(configPath) // nolint

	err := Configure("test_config.json")
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}
	if Config == nil || Config.LogLevel != "warn" {
		t.Errorf("Config not set correctly: %+v", Config)
	}
}

func TestGetConfigFile_Error(t *testing.T) {
	tmpHome := t.TempDir()
	oldUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) { return tmpHome, nil }
	defer func() { userHomeDir = oldUserHomeDir }()

	// Ensure no local config file exists
	os.Remove(filepath.Join("config", "test_config.json")) // nolint

	cfg, err := getConfigFile("test_config.json")
	if err == nil || cfg != nil {
		t.Error("Expected error when config files are missing")
	}
}
