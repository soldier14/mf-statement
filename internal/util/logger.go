package util

import (
	"io"
	"log/slog"
	"os"
)

// Logger wraps slog.Logger for our application
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new logger with slog
func NewLogger(level slog.Level, output io.Writer) *Logger {
	if output == nil {
		output = os.Stderr
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(output, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

// NewDefaultLogger creates a logger with default settings
func NewDefaultLogger() *Logger {
	return NewLogger(slog.LevelInfo, os.Stderr)
}

// NewDebugLogger creates a debug-level logger
func NewDebugLogger() *Logger {
	return NewLogger(slog.LevelDebug, os.Stderr)
}

// NewQuietLogger creates a logger that only shows errors
func NewQuietLogger() *Logger {
	return NewLogger(slog.LevelError, os.Stderr)
}

// WithFields creates a new logger with additional context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}

	newLogger := l.Logger.With(args...)
	return &Logger{Logger: newLogger}
}

// WithField creates a new logger with a single field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.Logger.With(key, value)
	return &Logger{Logger: newLogger}
}
