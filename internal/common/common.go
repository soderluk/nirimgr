// Package common contains commonly used functions.
package common

import (
	"log/slog"
	"os"
	"reflect"

	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/models"
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

// SetUintField is a helper function to set a uint64 field dynamically if present.
//
// Note that if the field already has a value, we don't want to override it.
func SetUintField(field reflect.Value, fieldName string, val any) {
	f := field.FieldByName(fieldName)
	slog.Debug("SetUintField", "field", field, "fieldName", fieldName, "value", val)
	if f.IsValid() && f.CanSet() {
		switch f.Kind() {
		case reflect.Uint8:
			if f.Uint() == 0 {
				f.SetUint(uint64(val.(uint8)))
			}
		case reflect.Uint64:
			if f.Uint() == 0 {
				f.SetUint(uint64(val.(uint64)))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if f.Int() == 0 {
				f.SetInt(int64(val.(int64)))
			}
		}
	}
}

// SetStringField is a helper function to set a string field dynamically if present.
func SetStringField(field reflect.Value, fieldName string, val string) {
	f := field.FieldByName(fieldName)
	slog.Debug("SetStringField", "field", field, "fieldName", fieldName, "value", val)
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
		if f.String() == "" {
			f.SetString(val)
		}
	}
}

// FilterWindows returns a slice of window models depending on the filtering function.
func FilterWindows(windows []*models.Window, f func(*models.Window) bool) []*models.Window {
	w := make([]*models.Window, 0)
	for _, e := range windows {
		if f(e) {
			w = append(w, e)
		}
	}

	return w
}

// FilterWindowsChain is a chainable function, you can use .First() to get the first window in the slice.
func FilterWindowsChain(windows []*models.Window, f func(*models.Window) bool) models.WindowSlice {
	return models.WindowSlice{Windows: FilterWindows(windows, f)}
}

// FilterWorkspaces returns a slice of workspace models depending on the filtering function.
func FilterWorkspaces(workspaces []*models.Workspace, f func(*models.Workspace) bool) []*models.Workspace {
	w := make([]*models.Workspace, 0)
	for _, e := range workspaces {
		if f(e) {
			w = append(w, e)
		}
	}

	return w
}

// FilterWorkspacesChain is a chainable function, you can use .First() to get the first workspace in the slice.
func FilterWorkspacesChain(workspaces []*models.Workspace, f func(*models.Workspace) bool) models.WorkspaceSlice {
	return models.WorkspaceSlice{Workspaces: FilterWorkspaces(workspaces, f)}
}

// FilterOutputs returns a slice of output models depending on the filtering function.
func FilterOutputs(outputs []*models.Output, f func(*models.Output) bool) []*models.Output {
	o := make([]*models.Output, 0)

	for _, e := range outputs {
		if f(e) {
			o = append(o, e)
		}
	}

	return o
}

// FilterOutputsChain is a chainable function, you can use .First() to get the first output in the slice.
func FilterOutputsChain(outputs []*models.Output, f func(*models.Output) bool) models.OutputSlice {
	return models.OutputSlice{Outputs: FilterOutputs(outputs, f)}
}
