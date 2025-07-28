package common

import (
	"log/slog"
	"testing"

	"github.com/soderluk/nirimgr/models"
	"github.com/stretchr/testify/assert"
)

func TestRepr(t *testing.T) {
	config := models.Config{}
	r := Repr(config)
	assert.Equal(t, "Config", r)
	r = Repr(nil)
	assert.Equal(t, "", r)
	window := &models.Window{ID: 1}
	r = Repr(window)
	assert.Equal(t, "Window", r)
}

func TestLogLevel(t *testing.T) {
	logLevels := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
		"FOO":   slog.LevelDebug,
	}
	for logLevel, expected := range logLevels {
		l := parseLogLevel(logLevel)
		assert.Equal(t, expected, l)
	}
}
