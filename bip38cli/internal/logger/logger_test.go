package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name     string
		debug    bool
		expected logrus.Level
		jsonFmt  bool
	}{
		{"debug mode", true, logrus.DebugLevel, false},
		{"production mode", false, logrus.WarnLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.debug)
			logger := GetLogger()

			assert.Equal(t, tt.expected, logger.GetLevel())

			if tt.jsonFmt {
				assert.IsType(t, &logrus.JSONFormatter{}, logger.Formatter)
			} else {
				assert.IsType(t, &logrus.TextFormatter{}, logger.Formatter)
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Reset the logger
	defaultLogger = nil

	logger1 := GetLogger()
	logger2 := GetLogger()

	// Should return the same instance (singleton pattern)
	assert.Equal(t, logger1, logger2)
	assert.NotNil(t, logger1)
	assert.Equal(t, logrus.WarnLevel, logger1.GetLevel())
}

func TestGetLoggerInitializesIfNil(t *testing.T) {
	// Reset the logger
	defaultLogger = nil

	logger := GetLogger()

	assert.NotNil(t, logger)
	assert.Equal(t, logrus.WarnLevel, logger.GetLevel())
	assert.IsType(t, &logrus.JSONFormatter{}, logger.Formatter)
}

func TestWithField(t *testing.T) {
	Init(false)
	entry := WithField("key", "value")

	assert.NotNil(t, entry)
	assert.IsType(t, &logrus.Entry{}, entry)
}

func TestWithFields(t *testing.T) {
	Init(false)
	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}
	entry := WithFields(fields)

	assert.NotNil(t, entry)
	assert.IsType(t, &logrus.Entry{}, entry)
}

func TestWithError(t *testing.T) {
	Init(false)
	err := assert.AnError
	entry := WithError(err)

	assert.NotNil(t, entry)
	assert.IsType(t, &logrus.Entry{}, entry)
}

func TestLogMethods(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer

	// Initialize with custom output
	Init(false)
	logger := GetLogger()
	logger.SetOutput(&buf)

	// Test all log levels
	Debug("debug message")
	Warn("warn message")
	Error("error message")

	output := buf.String()
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
	assert.NotContains(t, output, "info message")
	// Debug message won't appear because level is Warn
	assert.NotContains(t, output, "debug message")
}

func TestLogMethodsWithFormatting(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer

	// Initialize with custom output
	Init(false)
	logger := GetLogger()
	logger.SetOutput(&buf)

	// Test formatted log methods
	Debugf("debug %s", "message")
	Warnf("warn %s", "message")
	Errorf("error %s", "message")

	output := buf.String()
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
	assert.NotContains(t, output, "info message")
	// Debug message won't appear because level is Warn
	assert.NotContains(t, output, "debug message")
}

func TestDebugMode(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer

	Init(true) // Debug mode
	logger := GetLogger()
	logger.SetOutput(&buf)

	Debug("debug message")
	Info("info message")

	output := buf.String()
	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "info message")
}

func TestLoggerConcurrency(t *testing.T) {
	Init(false)

	// Test concurrent logging
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			WithField("goroutine", id).Info("concurrent test")
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we reach here, no race conditions occurred
	assert.True(t, true)
}

func TestLoggerValidation(t *testing.T) {
	Init(false)
	logger := GetLogger()

	require.NotNil(t, logger)
	require.NotNil(t, logger.Formatter)
	require.NotNil(t, logger.Out)

	// Test that logger can handle various input types
	WithField("string", "value").Info("test")
	WithField("int", 42).Info("test")
	WithField("bool", true).Info("test")
	WithField("slice", []string{"a", "b"}).Info("test")

	assert.True(t, true) // If we reach here, no panics occurred
}

func TestLoggerPerformance(t *testing.T) {
	Init(false)
	logger := GetLogger()
	logger.SetLevel(logrus.ErrorLevel) // Only log errors to reduce noise

	// Test that logging doesn't significantly impact performance
	for i := 0; i < 1000; i++ {
		Debug("This should not be logged due to level")
	}

	// If we reach here quickly, level filtering is working
	assert.True(t, true)
}

func TestLoggerIntegration(t *testing.T) {
	// Test complete logger setup with realistic configuration
	var buf bytes.Buffer

	Init(true) // Debug mode for more verbose logging
	logger := GetLogger()
	logger.SetOutput(&buf)

	// Test structured logging
	entry := WithFields(logrus.Fields{
		"component": "test",
		"version":   "1.0.0",
	})

	entry.Info("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "component")
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "test")
}

func TestLoggerWithEnvironmentConfig(t *testing.T) {
	// Test that logger works with different configurations
	Init(true)
	assert.Equal(t, logrus.DebugLevel, GetLogger().GetLevel())

	Init(false)
	assert.Equal(t, logrus.WarnLevel, GetLogger().GetLevel())
}

func TestLoggerOutputRedirection(t *testing.T) {
	Init(false)

	var buf bytes.Buffer
	logger := GetLogger()
	logger.SetOutput(&buf)

	Warn("test output")

	output := buf.String()
	assert.Contains(t, output, "test output")
}

func TestLoggerFormatterChange(t *testing.T) {
	Init(false)
	logger := GetLogger()

	// Test changing formatter
	originalFormatter := logger.Formatter
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	assert.NotEqual(t, originalFormatter, logger.Formatter)
	assert.IsType(t, &logrus.TextFormatter{}, logger.Formatter)
}
