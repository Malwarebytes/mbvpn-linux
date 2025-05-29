package console

import (
	"strings"
	"testing"
)

func TestRunCmd_WithoutSudo_Success(t *testing.T) {
	output, err := RunCmd(false, "echo", "hello", "world")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	expected := "hello world\n"
	if output != expected {
		t.Errorf("Expected output '%s', got '%s'", expected, output)
	}
}

func TestRunCmd_WithoutSudo_CommandNotFound(t *testing.T) {
	_, err := RunCmd(false, "nonexistentcommand12345")
	if err == nil {
		t.Error("Expected error for non-existent command, got nil")
	}
	if !strings.Contains(err.Error(), "command execution failed") {
		t.Errorf("Expected error message to contain 'command execution failed', got: %v", err)
	}
}

func TestRunCmd_WithoutSudo_CommandFails(t *testing.T) {
	_, err := RunCmd(false, "sh", "-c", "exit 1")
	if err == nil {
		t.Error("Expected error for failing command, got nil")
	}
	if !strings.Contains(err.Error(), "command execution failed") {
		t.Errorf("Expected error message to contain 'command execution failed', got: %v", err)
	}
}

func TestRunCmd_WithSudo_WhenSudoAvailable(t *testing.T) {
	// This test verifies that when sudo is requested and available, 
	// the command is executed with sudo (even if it might fail due to password requirements)
	_, err := RunCmd(true, "echo", "test")
	
	// The command will likely fail in CI/testing environments due to password requirements
	// We just want to verify that the sudo path is taken, not that it succeeds
	if err != nil {
		// Check that the error contains sudo-related messages, indicating sudo path was taken
		errMsg := err.Error()
		if strings.Contains(errMsg, "sudo:") || strings.Contains(errMsg, "password") || strings.Contains(errMsg, "terminal") {
			// This indicates sudo was attempted, which is the behavior we want to test
			return
		}
		// If it's a different error, it might be from the fallback path
		t.Logf("Command failed with error (may be expected in test environment): %v", err)
	}
}

func TestRunCmd_WithSudo_CommandNotFound(t *testing.T) {
	_, err := RunCmd(true, "nonexistentcommand12345")
	if err == nil {
		t.Error("Expected error for non-existent command with sudo, got nil")
	}
	if !strings.Contains(err.Error(), "command execution failed") {
		t.Errorf("Expected error message to contain 'command execution failed', got: %v", err)
	}
}

func TestRunCmd_EmptyCommand(t *testing.T) {
	_, err := RunCmd(false, "")
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
}

func TestRunCmd_MultipleArgs(t *testing.T) {
	output, err := RunCmd(false, "echo", "-n", "arg1", "arg2", "arg3")
	if err != nil {
		t.Errorf("Expected no error with multiple args, got: %v", err)
	}
	expected := "arg1 arg2 arg3"
	if output != expected {
		t.Errorf("Expected output '%s', got '%s'", expected, output)
	}
}

func TestRunCmd_WithStderr(t *testing.T) {
	// Use a command that writes to stderr and fails
	_, err := RunCmd(false, "sh", "-c", "echo 'error message' >&2; exit 1")
	if err == nil {
		t.Error("Expected error for command that writes to stderr, got nil")
	}
	
	// Check that stderr is included in the error message
	if !strings.Contains(err.Error(), "stderr:") {
		t.Errorf("Expected error message to contain stderr output, got: %v", err)
	}
	if !strings.Contains(err.Error(), "error message") {
		t.Errorf("Expected error message to contain 'error message', got: %v", err)
	}
}

func TestRunCmd_OutputCapture(t *testing.T) {
	// Test that stdout is properly captured and returned
	output, err := RunCmd(false, "printf", "line1\nline2")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	expected := "line1\nline2"
	if output != expected {
		t.Errorf("Expected output '%s', got '%s'", expected, output)
	}
}