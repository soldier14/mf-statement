package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger provides structured logging
type Logger struct {
	level  LogLevel
	output io.Writer
	logger *log.Logger
}

// NewLogger creates a new logger
func NewLogger(level LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stderr
	}

	return &Logger{
		level:  level,
		output: output,
		logger: log.New(output, "", 0),
	}
}

// NewDefaultLogger creates a logger with default settings
func NewDefaultLogger() *Logger {
	return NewLogger(LogLevelInfo, os.Stderr)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.log("DEBUG", format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.log("INFO", format, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.log("WARN", format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.log("ERROR", format, args...)
	}
}

// WithFields creates a new logger with additional context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// For simplicity, we'll just return the same logger
	// In a real implementation, you might want to use structured logging
	return l
}

func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] %s %s", timestamp, level, message)
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
	l.logger = log.New(output, "", 0)
}
