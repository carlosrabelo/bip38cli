package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// defaultLogger is the global logger instance used throughout the application
	defaultLogger *logrus.Logger
)

// Init initializes the default logger with the given configuration.
// If debug is true, logs are formatted as text with timestamps and debug level is enabled.
// If debug is false, logs are formatted as JSON and info level is used.
func Init(debug bool) {
	defaultLogger = logrus.New()

	// Set output to stdout
	defaultLogger.SetOutput(os.Stdout)

	// Set log level based on debug flag
	if debug {
		defaultLogger.SetLevel(logrus.DebugLevel)
		defaultLogger.Formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
	} else {
		defaultLogger.SetLevel(logrus.WarnLevel)
		defaultLogger.Formatter = &logrus.JSONFormatter{}
	}
}

// GetLogger returns the default logger instance.
// If the logger hasn't been initialized, it will be initialized with default settings.
func GetLogger() *logrus.Logger {
	if defaultLogger == nil {
		// Initialize with default settings if not already initialized
		Init(false)
	}
	return defaultLogger
}

// WithField creates a logger entry with a single field
func WithField(key string, value any) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields creates a logger entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithError creates a logger entry with an error field
func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}

// Debug logs a debug message
func Debug(args ...any) {
	GetLogger().Debug(args...)
}

// Info logs an info message
func Info(args ...any) {
	GetLogger().Info(args...)
}

// Warn logs a warning message
func Warn(args ...any) {
	GetLogger().Warn(args...)
}

// Error logs an error message
func Error(args ...any) {
	GetLogger().Error(args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...any) {
	GetLogger().Fatal(args...)
}

// Debugf logs a debug message with formatting
func Debugf(format string, args ...any) {
	GetLogger().Debugf(format, args...)
}

// Infof logs an info message with formatting
func Infof(format string, args ...any) {
	GetLogger().Infof(format, args...)
}

// Warnf logs a warning message with formatting
func Warnf(format string, args ...any) {
	GetLogger().Warnf(format, args...)
}

// Errorf logs an error message with formatting
func Errorf(format string, args ...any) {
	GetLogger().Errorf(format, args...)
}

// Fatalf logs a fatal message with formatting and exits
func Fatalf(format string, args ...any) {
	GetLogger().Fatalf(format, args...)
}
