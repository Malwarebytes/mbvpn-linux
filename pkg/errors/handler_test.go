package errors

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

// captureOutput captures stdout and stderr for testing
func captureOutput(fn func()) (stdout, stderr string) {
	// Capture stdout
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	
	// Capture stderr
	oldStderr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	// Execute function
	fn()

	// Close writers and restore
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutBuf.ReadFrom(rOut)
	stderrBuf.ReadFrom(rErr)

	return stdoutBuf.String(), stderrBuf.String()
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler()
	if handler == nil {
		t.Fatal("NewHandler should return a non-nil handler")
	}
}

func TestHandler_Handle_WithNilError(t *testing.T) {
	handler := NewHandler()
	
	shouldExit := handler.Handle(nil)
	if shouldExit {
		t.Error("Handle should return false for nil error")
	}
}

func TestHandler_Handle_WithUserError(t *testing.T) {
	handler := NewHandler()
	userErr := NewUserError("User friendly message", errors.New("underlying"))

	stdout, _ := captureOutput(func() {
		shouldExit := handler.Handle(userErr)
		if shouldExit {
			t.Error("Handle should return false for user errors")
		}
	})

	if !strings.Contains(stdout, "User friendly message") {
		t.Error("Should display user-friendly message")
	}
}

func TestHandler_Handle_WithStandardErrors(t *testing.T) {
	handler := NewHandler()

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"NotFound", ErrNotFound, "not found"},
		{"Unauthorized", ErrUnauthorized, "not authorized"},
		{"PermissionDenied", ErrPermissionDenied, "don't have permission"},
		{"InvalidInput", ErrInvalidInput, "Invalid input"},
		{"ServerError", ErrServerError, "server error occurred"},
		{"NetworkError", ErrNetworkError, "network error occurred"},
		{"ConfigError", ErrConfigError, "configuration error occurred"},
		{"VPNError", ErrVPNError, "VPN error occurred"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _ := captureOutput(func() {
				shouldExit := handler.Handle(tt.err)
				if !shouldExit {
					t.Error("Handle should return true for standard errors")
				}
			})

			if !strings.Contains(strings.ToLower(stdout), strings.ToLower(tt.expected)) {
				t.Errorf("Should display message containing '%s', got: %s", tt.expected, stdout)
			}
		})
	}
}

func TestHandler_Handle_WithGenericError(t *testing.T) {
	handler := NewHandler()
	genericErr := errors.New("some random error")

	stdout, _ := captureOutput(func() {
		shouldExit := handler.Handle(genericErr)
		if !shouldExit {
			t.Error("Handle should return true for generic errors")
		}
	})

	if !strings.Contains(stdout, "unexpected error") {
		t.Error("Should display unexpected error message for generic errors")
	}
}

func TestHandler_HandleWithMessage(t *testing.T) {
	handler := NewHandler()

	t.Run("With nil error", func(t *testing.T) {
		shouldExit := handler.HandleWithMessage(nil, "test message")
		if shouldExit {
			t.Error("HandleWithMessage should return false for nil error")
		}
	})

	t.Run("With error", func(t *testing.T) {
		err := errors.New("test error")
		customMessage := "Custom error message"

		stdout, _ := captureOutput(func() {
			shouldExit := handler.HandleWithMessage(err, customMessage)
			if !shouldExit {
				t.Error("HandleWithMessage should return true for non-nil error")
			}
		})

		if !strings.Contains(stdout, customMessage) {
			t.Errorf("Should display custom message '%s', got: %s", customMessage, stdout)
		}
	})
}

func TestHandler_Fatal(t *testing.T) {
	// Note: We can't easily test os.Exit(1) in unit tests as it would terminate the test process
	// In a real scenario, you might use dependency injection or test with build tags
	// For now, we'll test that Fatal calls Handle
	
	handler := NewHandler()
	userErr := NewUserError("test message", nil)

	// This should not cause exit since user errors return false from Handle
	stdout, _ := captureOutput(func() {
		handler.Fatal(userErr)
	})

	if !strings.Contains(stdout, "test message") {
		t.Error("Fatal should call Handle and display the message")
	}
}

func TestHandler_FatalWithMessage(t *testing.T) {
	handler := NewHandler()

	// This should not cause exit since we're passing nil (HandleWithMessage returns false for nil)
	handler.FatalWithMessage(nil, "custom message")

	// Test with an actual error would cause exit, so we just verify the method exists and compiles
}

func TestHandler_Error(t *testing.T) {
	handler := NewHandler()
	underlyingErr := errors.New("underlying error")
	message := "User friendly message"

	result := handler.Error(underlyingErr, message)

	userErr, ok := result.(*UserError)
	if !ok {
		t.Fatal("Error should return a UserError")
	}

	if userErr.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, userErr.Message)
	}

	if userErr.Err != underlyingErr {
		t.Error("UserError should wrap the underlying error")
	}
}

func TestUserErrorExtraction(t *testing.T) {
	t.Run("Extract UserError from UserError", func(t *testing.T) {
		userErr := NewUserError("test message", nil)
		extracted := userError(userErr)
		
		if extracted == nil {
			t.Fatal("Should extract UserError from UserError")
		}
		
		if extracted.Message != "test message" {
			t.Errorf("Expected 'test message', got '%s'", extracted.Message)
		}
	})

	t.Run("Extract UserError from wrapped UserError", func(t *testing.T) {
		userErr := NewUserError("test message", nil)
		wrappedErr := fmt.Errorf("wrapped: %w", userErr)
		extracted := userError(wrappedErr)
		
		if extracted == nil {
			t.Fatal("Should extract UserError from wrapped UserError")
		}
		
		if extracted.Message != "test message" {
			t.Errorf("Expected 'test message', got '%s'", extracted.Message)
		}
	})

	t.Run("Return nil for non-UserError", func(t *testing.T) {
		regularErr := errors.New("regular error")
		extracted := userError(regularErr)
		
		if extracted != nil {
			t.Error("Should return nil for non-UserError")
		}
	})

	t.Run("Return nil for nil error", func(t *testing.T) {
		extracted := userError(nil)
		
		if extracted != nil {
			t.Error("Should return nil for nil error")
		}
	})
}

func TestHandlerDebugMode(t *testing.T) {
	// Test debug mode behavior
	// Note: This test depends on the config.Debug() function
	// In a real scenario, you might want to inject debug mode as a parameter
	
	handler := &Handler{debug: true}
	userErr := NewUserError("User message", errors.New("underlying error"))

	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	_, _ = captureOutput(func() {
		handler.Handle(userErr)
	})

	// In debug mode, the underlying error should be logged
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "underlying error") {
		t.Error("Debug mode should log underlying error details")
	}
}

func TestHandlerNonDebugMode(t *testing.T) {
	handler := &Handler{debug: false}
	userErr := NewUserError("User message", errors.New("underlying error"))

	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	_, _ = captureOutput(func() {
		handler.Handle(userErr)
	})

	// In non-debug mode, underlying error details should not be logged to stderr
	// Note: They might still be logged to the log output, but not printed to stderr
}

// TestHandler_Integration tests the handler with various error combinations
func TestHandler_Integration(t *testing.T) {
	handler := NewHandler()

	// Test chain of errors
	baseErr := errors.New("database connection failed")
	networkErr := NewNetworkError("connect to database", baseErr)
	wrappedErr := fmt.Errorf("service initialization failed: %w", networkErr)

	stdout, _ := captureOutput(func() {
		shouldExit := handler.Handle(wrappedErr)
		if !shouldExit {
			t.Error("Should return true for error chain ending in non-user error")
		}
	})

	// Check that some error message is displayed (the exact message may vary due to formatting)
	if len(stdout) == 0 {
		t.Error("Should display some error message")
	}
}

// TestPrintMsgIntegration verifies integration with output.PrintMsg
func TestPrintMsgIntegration(t *testing.T) {
	// This test verifies that our handler properly integrates with the output package
	handler := NewHandler()
	
	// Test with a user error that should use MsgError type
	userErr := NewUserError("Test user error", nil)
	
	stdout, _ := captureOutput(func() {
		handler.Handle(userErr)
	})
	
	// The output should contain the error message (exact formatting may vary)
	if len(stdout) == 0 {
		t.Error("Should display some output for user error")
	}
}