package config

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempConfigFile(t *testing.T, content string) (string, func()) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test_config.json")
	err := os.WriteFile(file, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}
	return file, func() {
		err := os.Remove(file)
		if err != nil {
			t.Logf("Could not remove file %v", file)
		}
	}
}

func TestNewConfig_LocalFile(t *testing.T) {
	configContent := `{"logLevel":"debug"}`
	file, cleanup := createTempConfigFile(t, configContent)
	defer cleanup()

	// Move file to expected local path
	if err := os.MkdirAll("config", 0o755); err != nil {
		t.Log("could not create directory 'config'")
	}
	if err := os.Rename(file, "config/test_config.json"); err != nil {
		t.Logf("could not rename file %v", file)
	}
	defer func() {
		if err := os.Remove("config/test_config.json"); err != nil {
			t.Log("could not remove config file")
		}
	}()

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
	file, cleanup := createTempConfigFile(t, configContent)
	defer cleanup()

	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "nirimgr")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Logf("could not create directory '%v'", configDir)
	}
	configPath := filepath.Join(configDir, "test_config.json")
	if err := os.Rename(file, configPath); err != nil {
		t.Logf("could not rename file '%v' to '%v'", file, configPath)
	}
	defer func() {
		if err := os.Remove(configPath); err != nil {
			t.Log("could not remove config path")
		}
	}()

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
	file, cleanup := createTempConfigFile(t, configContent)
	defer cleanup()
	if err := os.MkdirAll("config", 0o755); err != nil {
		t.Log("could not create directory 'config'")
	}
	if err := os.Rename(file, "config/test_config.json"); err != nil {
		t.Logf("could not rename file '%v'", file)
	}
	defer func() {
		if err := os.Remove("config/test_config.json"); err != nil {
			t.Log("could not remove test config file")
		}
	}()

	err := Configure("test_config.json")
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}
	if Config == nil || Config.LogLevel != "warn" {
		t.Errorf("Config not set correctly: %+v", Config)
	}
}

func TestGetConfigFile_Error(t *testing.T) {
	err := os.Remove("config/test_config.json")
	if err != nil {
		t.Log("could not remove test config.")
	}
	cfg, err := getConfigFile("test_config.json")
	if err == nil || cfg != nil {
		t.Error("Expected error when config files are missing")
	}
}
