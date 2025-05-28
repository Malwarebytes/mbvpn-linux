package console

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestRunCmd_WithoutSudo(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "successful echo command",
			command:     "echo",
			args:        []string{"hello", "world"},
			expectError: false,
		},
		{
			name:        "successful command with no args",
			command:     "echo",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "nonexistent command",
			command:     "nonexistent-command-xyz",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "command that fails",
			command:     "ls",
			args:        []string{"/nonexistent-directory-xyz"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RunCmd(false, tt.command, tt.args...)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				
				// For echo commands, verify output contains expected text
				if tt.command == "echo" && len(tt.args) > 0 {
					expectedOutput := strings.Join(tt.args, " ")
					if !strings.Contains(output, expectedOutput) {
						t.Errorf("Expected output to contain '%s', got '%s'", expectedOutput, output)
					}
				}
			}
		})
	}
}

func TestRunCmd_WithSudo(t *testing.T) {
	// Skip sudo tests on Windows as sudo is not available
	if runtime.GOOS == "windows" {
		t.Skip("Skipping sudo tests on Windows")
	}

	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "sudo with echo command",
			command:     "echo",
			args:        []string{"hello", "sudo"},
			expectError: false,
		},
		{
			name:        "sudo with nonexistent command",
			command:     "nonexistent-command-xyz",
			args:        []string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := RunCmd(true, tt.command, tt.args...)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				
				// For echo commands, verify output contains expected text
				if tt.command == "echo" && len(tt.args) > 0 {
					expectedOutput := strings.Join(tt.args, " ")
					if !strings.Contains(output, expectedOutput) {
						t.Errorf("Expected output to contain '%s', got '%s'", expectedOutput, output)
					}
				}
			}
		})
	}
}

func TestRunCmd_SudoFallback(t *testing.T) {
	// This test verifies that when sudo is not available, the command runs without sudo
	// We temporarily modify PATH to exclude sudo
	if runtime.GOOS == "windows" {
		t.Skip("Skipping sudo fallback test on Windows")
	}

	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Set PATH to empty to simulate sudo not being available
	os.Setenv("PATH", "")

	output, err := RunCmd(true, "echo", "test")
	if err != nil {
		t.Errorf("Expected command to run without sudo when sudo is not available, got error: %v", err)
	}

	if !strings.Contains(output, "test") {
		t.Errorf("Expected output to contain 'test', got '%s'", output)
	}
}

func TestRunCmd_StdoutStderr(t *testing.T) {
	// Test that both stdout and stderr are captured in error messages
	// Use a command that produces stderr
	_, err := RunCmd(false, "ls", "/nonexistent-directory-xyz")
	
	if err == nil {
		t.Error("Expected error for nonexistent directory")
		return
	}

	errorMsg := err.Error()
	
	// Check that error message contains the expected format
	if !strings.Contains(errorMsg, "command execution failed") {
		t.Errorf("Expected error message to contain 'command execution failed', got: %s", errorMsg)
	}
	
	if !strings.Contains(errorMsg, "stderr:") {
		t.Errorf("Expected error message to contain 'stderr:', got: %s", errorMsg)
	}
	
	if !strings.Contains(errorMsg, "stdout:") {
		t.Errorf("Expected error message to contain 'stdout:', got: %s", errorMsg)
	}
}

func TestRunCmd_EmptyArgs(t *testing.T) {
	// Test that command works with no arguments
	output, err := RunCmd(false, "echo")
	if err != nil {
		t.Errorf("Expected no error for echo with no args, got: %v", err)
	}
	
	// echo with no args should produce a newline
	if output != "\n" {
		t.Errorf("Expected output to be a newline, got: '%s'", output)
	}
}

func TestRunCmd_MultipleArgs(t *testing.T) {
	// Test command with multiple arguments
	output, err := RunCmd(false, "echo", "arg1", "arg2", "arg3")
	if err != nil {
		t.Errorf("Expected no error for echo with multiple args, got: %v", err)
	}
	
	expectedOutput := "arg1 arg2 arg3"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain '%s', got '%s'", expectedOutput, output)
	}
}

// TestRunCmd_Integration tests the function with realistic commands
func TestRunCmd_Integration(t *testing.T) {
	// Test with 'which' command to find a common binary
	output, err := RunCmd(false, "which", "echo")
	if err != nil {
		t.Errorf("Expected 'which echo' to succeed, got error: %v", err)
	}
	
	if !strings.Contains(output, "echo") {
		t.Errorf("Expected output to contain 'echo', got: '%s'", output)
	}
	
	// Test with 'whoami' command
	output, err = RunCmd(false, "whoami")
	if err != nil {
		t.Errorf("Expected 'whoami' to succeed, got error: %v", err)
	}
	
	if strings.TrimSpace(output) == "" {
		t.Error("Expected 'whoami' to produce non-empty output")
	}
}