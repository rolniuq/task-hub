package logger

import (
	"log/slog"

	"go.uber.org/fx"
)

var LoggerModule = fx.Module(
	"logger",
	fx.Provide(NewLogger),
)

type Logger struct {
	slog *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		slog: slog.Default(),
	}
}

func (l *Logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}
