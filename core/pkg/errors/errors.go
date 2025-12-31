// Package errors defines the error types and handling for the application.
package errors

import (
	"fmt"
)

// ErrorType represents different categories of errors in the application.
// It provides a way to categorize and handle different types of errors consistently.
type ErrorType string

const (
	// ValidationError occurs when input validation fails
	ValidationError ErrorType = "validation"

	// CryptoError occurs during cryptographic operations
	CryptoError ErrorType = "crypto"

	// ConfigError occurs during configuration loading
	ConfigError ErrorType = "config"

	// InputError occurs when reading user input
	InputError ErrorType = "input"

	// SystemError occurs for system-level failures
	SystemError ErrorType = "system"
)

// AppError represents a structured application error with type, message, cause, and context.
// It provides rich error information for better debugging and error handling.
type AppError struct {
	Type    ErrorType      // The category of error
	Message string         // Human-readable error message
	Cause   error          // The underlying error that caused this error
	Context map[string]any // Additional context information
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error with the given message and cause.
// Validation errors occur when input validation fails.
func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// NewCryptoError creates a new crypto error with the given message and cause.
// Crypto errors occur during cryptographic operations.
func NewCryptoError(message string, cause error) *AppError {
	return &AppError{
		Type:    CryptoError,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// NewConfigError creates a new config error with the given message and cause.
// Config errors occur during configuration loading or parsing.
func NewConfigError(message string, cause error) *AppError {
	return &AppError{
		Type:    ConfigError,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// NewInputError creates a new input error with the given message and cause.
// Input errors occur when reading user input or processing input data.
func NewInputError(message string, cause error) *AppError {
	return &AppError{
		Type:    InputError,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// NewSystemError creates a new system error with the given message and cause.
// System errors occur for system-level failures like file I/O or network issues.
func NewSystemError(message string, cause error) *AppError {
	return &AppError{
		Type:    SystemError,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to the error and returns the error for chaining.
// Context provides additional debugging information about when and where the error occurred.
func (e *AppError) WithContext(key string, value any) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// GetContext returns the error context
func (e *AppError) GetContext() map[string]any {
	return e.Context
}

// IsType checks if the error is of a specific type
func (e *AppError) IsType(errorType ErrorType) bool {
	return e.Type == errorType
}

// Helper functions to check error types

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.IsType(ValidationError)
	}
	return false
}

// IsCryptoError checks if an error is a crypto error
func IsCryptoError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.IsType(CryptoError)
	}
	return false
}

// IsConfigError checks if an error is a config error
func IsConfigError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.IsType(ConfigError)
	}
	return false
}

// IsInputError checks if an error is an input error
func IsInputError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.IsType(InputError)
	}
	return false
}

// IsSystemError checks if an error is a system error
func IsSystemError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.IsType(SystemError)
	}
	return false
}
