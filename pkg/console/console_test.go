package console

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test command allowlisting - only wg and wg-quick should be allowed
func TestRunCmd_DisallowedCommand(t *testing.T) {
	_, err := RunCmd("echo", "test")
	if err == nil {
		t.Error("Expected error for disallowed command 'echo', got nil")
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Errorf("Expected error about disallowed command, got: %v", err)
	}
}

func TestRunCmd_DisallowedCommandWithSudo(t *testing.T) {
	_, err := RunCmd("bash", "-c", "echo test")
	if err == nil {
		t.Error("Expected error for disallowed command 'bash', got nil")
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Errorf("Expected error about disallowed command, got: %v", err)
	}
}

// Test empty command validation
func TestRunCmd_EmptyCommand(t *testing.T) {
	_, err := RunCmd("")
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected error about empty command, got: %v", err)
	}
}

// Test wg command validation
func TestRunCmd_WgShow_Valid(t *testing.T) {
	// This may fail if wg is not installed, but should pass validation
	_, err := RunCmd("wg", "show")
	// We accept either success or execution failure (if wg not installed)
	// but NOT validation errors
	if err != nil && strings.Contains(err.Error(), "not allowed") {
		t.Errorf("wg show should pass validation, got: %v", err)
	}
}

func TestRunCmd_WgShow_NoArgs(t *testing.T) {
	// wg with no args should also be valid (defaults to show)
	_, err := RunCmd("wg")
	// We accept either success or execution failure (if wg not installed)
	// but NOT validation errors
	if err != nil && strings.Contains(err.Error(), "not allowed") {
		t.Errorf("wg with no args should pass validation, got: %v", err)
	}
}

func TestRunCmd_Wg_InvalidSubcommand(t *testing.T) {
	// wg with other subcommands should be rejected
	_, err := RunCmd("wg", "set", "wg0", "peer", "xxxxx")
	if err == nil {
		t.Error("Expected error for wg with unsupported subcommand, got nil")
	}
	if !strings.Contains(err.Error(), "only supports 'show'") {
		t.Errorf("Expected error about unsupported wg command, got: %v", err)
	}
}

// Test wg-quick command validation
func TestRunCmd_WgQuick_Valid(t *testing.T) {
	// Create a temporary config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")
	err := os.WriteFile(configPath, []byte("[Interface]\n"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// This may fail if wg-quick is not installed, but should pass validation
	_, err = RunCmd("wg-quick", "up", configPath)
	// We accept either success or execution failure (if wg-quick not installed)
	// but NOT validation errors about "not allowed"
	if err != nil && strings.Contains(err.Error(), "not allowed") {
		t.Errorf("wg-quick up should pass validation, got: %v", err)
	}
}

func TestRunCmd_WgQuick_Down_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")
	err := os.WriteFile(configPath, []byte("[Interface]\n"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err = RunCmd("wg-quick", "down", configPath)
	// We accept either success or execution failure (if wg-quick not installed)
	// but NOT validation errors
	if err != nil && strings.Contains(err.Error(), "not allowed") {
		t.Errorf("wg-quick down should pass validation, got: %v", err)
	}
}

func TestRunCmd_WgQuick_InvalidAction(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")

	_, err := RunCmd("wg-quick", "restart", configPath)
	if err == nil {
		t.Error("Expected error for invalid wg-quick action, got nil")
	}
	if !strings.Contains(err.Error(), "must be 'up' or 'down'") {
		t.Errorf("Expected error about invalid action, got: %v", err)
	}
}

func TestRunCmd_WgQuick_RelativePath(t *testing.T) {
	_, err := RunCmd("wg-quick", "up", "relative/path/test.conf")
	if err == nil {
		t.Error("Expected error for relative path, got nil")
	}
	if !strings.Contains(err.Error(), "must be absolute") {
		t.Errorf("Expected error about absolute path requirement, got: %v", err)
	}
}

func TestRunCmd_WgQuick_NonConfFile(t *testing.T) {
	_, err := RunCmd("wg-quick", "up", "/tmp/test.txt")
	if err == nil {
		t.Error("Expected error for non-.conf file, got nil")
	}
	if !strings.Contains(err.Error(), ".conf extension") {
		t.Errorf("Expected error about .conf extension, got: %v", err)
	}
}

func TestRunCmd_WgQuick_PathTraversal(t *testing.T) {
	_, err := RunCmd("wg-quick", "up", "/tmp/../etc/test.conf")
	if err == nil {
		t.Error("Expected error for path traversal attempt, got nil")
	}
	if !strings.Contains(err.Error(), "invalid config file path") {
		t.Errorf("Expected error about invalid path, got: %v", err)
	}
}

func TestRunCmd_WgQuick_MissingArgs(t *testing.T) {
	_, err := RunCmd("wg-quick", "up")
	if err == nil {
		t.Error("Expected error for missing config file argument, got nil")
	}
	if !strings.Contains(err.Error(), "requires action") {
		t.Errorf("Expected error about missing arguments, got: %v", err)
	}
}

// Test shell metacharacter injection prevention
func TestRunCmd_ShellMetacharacters(t *testing.T) {
	testCases := []struct {
		name    string
		command string
		args    []string
	}{
		{"semicolon", "wg", []string{"show;echo test"}},
		{"pipe", "wg", []string{"show|cat"}},
		{"redirect", "wg", []string{"show>/tmp/out"}},
		{"backtick", "wg", []string{"show`whoami`"}},
		{"dollar", "wg", []string{"show$(whoami)"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := RunCmd(tc.command, tc.args...)
			if err == nil {
				t.Errorf("Expected error for shell metacharacter in %s, got nil", tc.name)
			}
			// Should fail validation, not execution
			if strings.Contains(err.Error(), "command execution failed") {
				t.Errorf("Should fail validation before execution for %s, got: %v", tc.name, err)
			}
		})
	}
}

// Test sudo requirement
func TestRunCmd_SudoNotAvailable(t *testing.T) {
	// Temporarily modify PATH to hide sudo
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", oldPath)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")

	_, err := RunCmd("wg-quick", "up", configPath)
	if err == nil {
		t.Error("Expected error when sudo is not available, got nil")
	}
	if !strings.Contains(err.Error(), "sudo is required but not available") {
		t.Errorf("Expected error about sudo not available, got: %v", err)
	}
}

// Test that validation happens before execution
func TestRunCmd_ValidationBeforeExecution(t *testing.T) {
	// Use an invalid command - should fail validation, not execution
	_, err := RunCmd("ls", "/tmp")
	if err == nil {
		t.Error("Expected error for disallowed command, got nil")
	}
	// Should be validation error, not execution error
	if !strings.Contains(err.Error(), "not allowed") {
		t.Errorf("Expected validation error, got: %v", err)
	}
	if strings.Contains(err.Error(), "command execution failed") {
		t.Error("Command should not have been executed")
	}
}
