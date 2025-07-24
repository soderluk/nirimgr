// Package config handles the configuration of nirimgr.
//
// See an example configuration in the README.md.
package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/soderluk/nirimgr/models"
)

var (
	// Version contains the current version of nirimgr
	Version string = "dev"
	// BuildDate is the date when nirimgr was built
	BuildDate string = "unknown"
	// CommitSHA contains the build commit SHA hash.
	CommitSHA string = "none"
	// Config contains all configurations
	Config *models.Config
)

// getConfigFile returns the config.json file.
//
// We first try locally, if we're e.g. running nirimgr with go run main.go,
// if that's not found, try ~/.config/nirimgr/config.json.
func getConfigFile(filename string) (*os.File, error) {
	if !strings.Contains(filename, "config.json") {
		slog.Error("Invalid configuration name", "got", filename, "want", "*config.json")
		return nil, fmt.Errorf("invalid configuration filename")
	}
	f, err := os.Open("config/" + filename) // #nosec G304
	if err != nil {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return nil, homeErr
		}
		configPath := filepath.Join(homeDir, ".config", "nirimgr", filename)
		f, err = os.Open(configPath) // #nosec G304
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

// newConfig configures the application.
//
// Returns the decoded data from the specified config file in the config struct.
func newConfig(filename string) (*models.Config, error) {
	f, err := getConfigFile(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Could not close config file", "error", err.Error())
		}
	}()

	var c *models.Config
	err = json.NewDecoder(f).Decode(&c)

	slog.Debug("Configured", "config", c)
	return c, err
}

// Configure reads the configuration file (json) to get the configuration.
//
// Sets the global Config, so it can be accessed from anywhere.
// See example configuration in the README.md.
func Configure(filename string) error {
	cfg, err := newConfig(filename)
	if err != nil {
		slog.Error("Could not read configuration from file.", "error", err.Error())
		return err
	}
	Config = cfg
	return nil
}
