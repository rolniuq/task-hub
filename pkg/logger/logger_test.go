package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.slog)
}

func TestLogger_Error(t *testing.T) {
	logger := NewLogger()
	assert.NotPanics(t, func() {
		logger.Error("test error message", "key", "value")
	})
}

func TestLogger_Info(t *testing.T) {
	logger := NewLogger()
	assert.NotPanics(t, func() {
		logger.Info("test info message", "key", "value")
	})
}

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger()
	assert.NotPanics(t, func() {
		logger.Debug("test debug message", "key", "value")
	})
}

func TestLoggerModule(t *testing.T) {
	assert.NotNil(t, LoggerModule)
}
