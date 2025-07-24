// Package common contains commonly used functions.
package common

import (
	"log/slog"
	"os"
	"reflect"

	"github.com/soderluk/nirimgr/config"
)

// SetupLogger sets up the logging for the application.
//
// The log level can be defined as the string "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"
// in the configuration file. E.g. "LogLevel": "INFO".
// Defaults to "DEBUG".
func SetupLogger() {
	logLevel := parseLogLevel(config.Config.LogLevel)
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(handler)

	slog.SetDefault(logger)
}

// Repr returns the name of the given model.
//
// This can be used to print out the model name.
func Repr(model any) string {
	if model == nil {
		return ""
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

// parseLogLevel parses the given log level string to slog log level.
func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
