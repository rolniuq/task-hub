package logger

import "log/slog"

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
