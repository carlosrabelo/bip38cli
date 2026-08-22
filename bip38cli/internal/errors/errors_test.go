package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	// Test error without cause
	err := NewValidationError("test message", nil)
	expected := "validation: test message"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test error with cause
	cause := errors.New("underlying error")
	err = NewCryptoError("crypto failed", cause)
	expected = "crypto: crypto failed (caused by: underlying error)"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewValidationError("test message", cause)

	if err.Unwrap() != cause {
		t.Error("Unwrap should return the underlying cause")
	}

	// Test without cause
	errNoCause := NewValidationError("test message", nil)
	if errNoCause.Unwrap() != nil {
		t.Error("Unwrap should return nil when there's no cause")
	}
}

func TestAppError_IsType(t *testing.T) {
	err := NewValidationError("test message", nil)

	if !err.IsType(ValidationError) {
		t.Error("Error should be of type ValidationError")
	}

	if err.IsType(CryptoError) {
		t.Error("Error should not be of type CryptoError")
	}
}

func TestAppError_WithContext(t *testing.T) {
	err := NewValidationError("test message", nil)
	err = err.WithContext("key", "value").WithContext("number", 42)

	context := err.GetContext()
	if context["key"] != "value" {
		t.Errorf("Expected context key 'key' to be 'value', got %v", context["key"])
	}

	if context["number"] != 42 {
		t.Errorf("Expected context key 'number' to be 42, got %v", context["number"])
	}
}

func TestNewValidationError(t *testing.T) {
	cause := errors.New("test cause")
	err := NewValidationError("validation failed", cause)

	if !err.IsType(ValidationError) {
		t.Error("Should create ValidationError")
	}

	if err.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", err.Message)
	}

	if err.Cause != cause {
		t.Error("Should set the cause correctly")
	}
}

func TestNewCryptoError(t *testing.T) {
	err := NewCryptoError("crypto failed", nil)

	if !err.IsType(CryptoError) {
		t.Error("Should create CryptoError")
	}

	if err.Message != "crypto failed" {
		t.Errorf("Expected message 'crypto failed', got '%s'", err.Message)
	}
}

func TestNewConfigError(t *testing.T) {
	err := NewConfigError("config invalid", nil)

	if !err.IsType(ConfigError) {
		t.Error("Should create ConfigError")
	}
}

func TestNewInputError(t *testing.T) {
	err := NewInputError("input invalid", nil)

	if !err.IsType(InputError) {
		t.Error("Should create InputError")
	}
}

func TestNewSystemError(t *testing.T) {
	err := NewSystemError("system error", nil)

	if !err.IsType(SystemError) {
		t.Error("Should create SystemError")
	}
}

func TestIsValidationError(t *testing.T) {
	err := NewValidationError("test", nil)
	if !IsValidationError(err) {
		t.Error("Should return true for ValidationError")
	}

	otherErr := NewCryptoError("test", nil)
	if IsValidationError(otherErr) {
		t.Error("Should return false for non-ValidationError")
	}

	plainErr := errors.New("plain error")
	if IsValidationError(plainErr) {
		t.Error("Should return false for plain error")
	}
}

func TestIsCryptoError(t *testing.T) {
	err := NewCryptoError("test", nil)
	if !IsCryptoError(err) {
		t.Error("Should return true for CryptoError")
	}

	otherErr := NewValidationError("test", nil)
	if IsCryptoError(otherErr) {
		t.Error("Should return false for non-CryptoError")
	}
}

func TestIsConfigError(t *testing.T) {
	err := NewConfigError("test", nil)
	if !IsConfigError(err) {
		t.Error("Should return true for ConfigError")
	}
}

func TestIsInputError(t *testing.T) {
	err := NewInputError("test", nil)
	if !IsInputError(err) {
		t.Error("Should return true for InputError")
	}
}

func TestIsSystemError(t *testing.T) {
	err := NewSystemError("test", nil)
	if !IsSystemError(err) {
		t.Error("Should return true for SystemError")
	}
}
