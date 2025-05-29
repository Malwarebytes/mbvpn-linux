package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestStandardErrors(t *testing.T) {
	// Test that all standard errors are defined and non-nil
	standardErrors := []error{
		ErrNotFound,
		ErrUnauthorized,
		ErrPermissionDenied,
		ErrInvalidInput,
		ErrServerError,
		ErrNetworkError,
		ErrConfigError,
		ErrVPNError,
	}

	for _, err := range standardErrors {
		if err == nil {
			t.Error("Standard error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Standard error should have a non-empty message")
		}
	}
}

func TestConfigError(t *testing.T) {
	t.Run("ConfigError with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("underlying error")
		configErr := NewConfigError("test operation", underlyingErr)

		expected := "config error during test operation: underlying error"
		if configErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, configErr.Error())
		}

		if configErr.Unwrap() != underlyingErr {
			t.Error("Unwrap should return the underlying error")
		}
	})

	t.Run("ConfigError with nil underlying error", func(t *testing.T) {
		configErr := NewConfigError("test operation", nil)

		expected := "config error during test operation"
		if configErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, configErr.Error())
		}

		if configErr.Unwrap() != nil {
			t.Error("Unwrap should return nil when underlying error is nil")
		}
	})

	t.Run("ConfigError struct creation", func(t *testing.T) {
		configErr := &ConfigError{Op: "direct creation", Err: fmt.Errorf("test error")}
		
		if configErr.Op != "direct creation" {
			t.Errorf("Expected operation 'direct creation', got '%s'", configErr.Op)
		}
		
		if configErr.Err.Error() != "test error" {
			t.Errorf("Expected error 'test error', got '%s'", configErr.Err.Error())
		}
	})
}

func TestAuthError(t *testing.T) {
	t.Run("AuthError with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("auth failed")
		authErr := NewAuthError("login", underlyingErr)

		expected := "auth error during login: auth failed"
		if authErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, authErr.Error())
		}

		if authErr.Unwrap() != underlyingErr {
			t.Error("Unwrap should return the underlying error")
		}
	})

	t.Run("AuthError with nil underlying error", func(t *testing.T) {
		authErr := NewAuthError("logout", nil)

		expected := "auth error during logout"
		if authErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, authErr.Error())
		}

		if authErr.Unwrap() != nil {
			t.Error("Unwrap should return nil when underlying error is nil")
		}
	})
}

func TestNetworkError(t *testing.T) {
	t.Run("NetworkError with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("connection timeout")
		networkErr := NewNetworkError("api call", underlyingErr)

		expected := "network error during api call: connection timeout"
		if networkErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, networkErr.Error())
		}

		if networkErr.Unwrap() != underlyingErr {
			t.Error("Unwrap should return the underlying error")
		}
	})

	t.Run("NetworkError with nil underlying error", func(t *testing.T) {
		networkErr := NewNetworkError("ping", nil)

		expected := "network error during ping"
		if networkErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, networkErr.Error())
		}

		if networkErr.Unwrap() != nil {
			t.Error("Unwrap should return nil when underlying error is nil")
		}
	})
}

func TestVPNError(t *testing.T) {
	t.Run("VPNError with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("wireguard failed")
		vpnErr := NewVPNError("connect", underlyingErr)

		expected := "vpn error during connect: wireguard failed"
		if vpnErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, vpnErr.Error())
		}

		if vpnErr.Unwrap() != underlyingErr {
			t.Error("Unwrap should return the underlying error")
		}
	})

	t.Run("VPNError with nil underlying error", func(t *testing.T) {
		vpnErr := NewVPNError("disconnect", nil)

		expected := "vpn error during disconnect"
		if vpnErr.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, vpnErr.Error())
		}

		if vpnErr.Unwrap() != nil {
			t.Error("Unwrap should return nil when underlying error is nil")
		}
	})
}

func TestUserError(t *testing.T) {
	t.Run("UserError with underlying error", func(t *testing.T) {
		underlyingErr := errors.New("detailed error")
		userErr := NewUserError("User friendly message", underlyingErr)

		if userErr.Error() != "User friendly message" {
			t.Errorf("Expected 'User friendly message', got '%s'", userErr.Error())
		}

		if userErr.Unwrap() != underlyingErr {
			t.Error("Unwrap should return the underlying error")
		}

		if userErr.Message != "User friendly message" {
			t.Errorf("Expected message 'User friendly message', got '%s'", userErr.Message)
		}

		if userErr.Err != underlyingErr {
			t.Error("Err field should contain the underlying error")
		}
	})

	t.Run("UserError with nil underlying error", func(t *testing.T) {
		userErr := NewUserError("Simple message", nil)

		if userErr.Error() != "Simple message" {
			t.Errorf("Expected 'Simple message', got '%s'", userErr.Error())
		}

		if userErr.Unwrap() != nil {
			t.Error("Unwrap should return nil when underlying error is nil")
		}
	})
}

func TestIsUserError(t *testing.T) {
	t.Run("Returns true for UserError", func(t *testing.T) {
		userErr := NewUserError("test message", nil)
		if !IsUserError(userErr) {
			t.Error("IsUserError should return true for UserError")
		}
	})

	t.Run("Returns false for non-UserError", func(t *testing.T) {
		regularErr := errors.New("regular error")
		if IsUserError(regularErr) {
			t.Error("IsUserError should return false for regular error")
		}
	})

	t.Run("Returns false for nil error", func(t *testing.T) {
		if IsUserError(nil) {
			t.Error("IsUserError should return false for nil error")
		}
	})

	t.Run("Returns true for wrapped UserError", func(t *testing.T) {
		userErr := NewUserError("test message", nil)
		wrappedErr := fmt.Errorf("wrapped: %w", userErr)
		if !IsUserError(wrappedErr) {
			t.Error("IsUserError should return true for wrapped UserError")
		}
	})
}

func TestGetUserErrorMessage(t *testing.T) {
	t.Run("Returns user message for UserError", func(t *testing.T) {
		userErr := NewUserError("User friendly message", errors.New("detailed"))
		message := GetUserErrorMessage(userErr)
		if message != "User friendly message" {
			t.Errorf("Expected 'User friendly message', got '%s'", message)
		}
	})

	t.Run("Returns default message for non-UserError", func(t *testing.T) {
		regularErr := errors.New("regular error")
		message := GetUserErrorMessage(regularErr)
		if message != "An unexpected error occurred" {
			t.Errorf("Expected 'An unexpected error occurred', got '%s'", message)
		}
	})

	t.Run("Returns default message for nil error", func(t *testing.T) {
		message := GetUserErrorMessage(nil)
		if message != "An unexpected error occurred" {
			t.Errorf("Expected 'An unexpected error occurred', got '%s'", message)
		}
	})

	t.Run("Returns user message for wrapped UserError", func(t *testing.T) {
		userErr := NewUserError("User message", nil)
		wrappedErr := fmt.Errorf("wrapped: %w", userErr)
		message := GetUserErrorMessage(wrappedErr)
		if message != "User message" {
			t.Errorf("Expected 'User message', got '%s'", message)
		}
	})
}

func TestErrorCheckFunctions(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		checkFn  func(error) bool
		expected bool
	}{
		{"IsNotFound with ErrNotFound", ErrNotFound, IsNotFound, true},
		{"IsNotFound with other error", ErrUnauthorized, IsNotFound, false},
		{"IsNotFound with wrapped ErrNotFound", fmt.Errorf("wrapped: %w", ErrNotFound), IsNotFound, true},
		{"IsUnauthorized with ErrUnauthorized", ErrUnauthorized, IsUnauthorized, true},
		{"IsUnauthorized with other error", ErrNotFound, IsUnauthorized, false},
		{"IsPermissionDenied with ErrPermissionDenied", ErrPermissionDenied, IsPermissionDenied, true},
		{"IsPermissionDenied with other error", ErrNotFound, IsPermissionDenied, false},
		{"IsInvalidInput with ErrInvalidInput", ErrInvalidInput, IsInvalidInput, true},
		{"IsInvalidInput with other error", ErrNotFound, IsInvalidInput, false},
		{"IsServerError with ErrServerError", ErrServerError, IsServerError, true},
		{"IsServerError with other error", ErrNotFound, IsServerError, false},
		{"IsNetworkError with ErrNetworkError", ErrNetworkError, IsNetworkError, true},
		{"IsNetworkError with other error", ErrNotFound, IsNetworkError, false},
		{"IsConfigError with ErrConfigError", ErrConfigError, IsConfigError, true},
		{"IsConfigError with other error", ErrNotFound, IsConfigError, false},
		{"IsVPNError with ErrVPNError", ErrVPNError, IsVPNError, true},
		{"IsVPNError with other error", ErrNotFound, IsVPNError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checkFn(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for %s", tt.expected, result, tt.name)
			}
		})
	}
}

func TestErrorCheckFunctionsWithNil(t *testing.T) {
	checkFunctions := []struct {
		name string
		fn   func(error) bool
	}{
		{"IsNotFound", IsNotFound},
		{"IsUnauthorized", IsUnauthorized},
		{"IsPermissionDenied", IsPermissionDenied},
		{"IsInvalidInput", IsInvalidInput},
		{"IsServerError", IsServerError},
		{"IsNetworkError", IsNetworkError},
		{"IsConfigError", IsConfigError},
		{"IsVPNError", IsVPNError},
	}

	for _, test := range checkFunctions {
		t.Run(test.name+" with nil", func(t *testing.T) {
			if test.fn(nil) {
				t.Errorf("%s should return false for nil error", test.name)
			}
		})
	}
}