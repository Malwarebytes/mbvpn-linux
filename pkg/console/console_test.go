package console

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test wg-quick with invalid action
func TestRunWgQuick_InvalidAction(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")

	_, err := RunWgQuick("restart", configPath)
	if err == nil {
		t.Error("Expected error for invalid wg-quick action, got nil")
	}
	if !strings.Contains(err.Error(), "must be 'up' or 'down'") {
		t.Errorf("Expected error about invalid action, got: %v", err)
	}
}

// Test wg command validation
func TestRunWg_ShowValid(t *testing.T) {
	// This may fail if wg is not installed, but should pass validation
	_, err := RunWg("show")
	// We accept either success or execution failure (if wg not installed)
	// but NOT validation errors
	if err != nil && strings.Contains(err.Error(), "only supports") {
		t.Errorf("wg show should pass validation, got: %v", err)
	}
}

func TestRunWg_NoArgs(t *testing.T) {
	// wg with no args should also be valid (defaults to show)
	_, err := RunWg()
	// We accept either success or execution failure (if wg not installed)
	// but NOT validation errors
	if err != nil && strings.Contains(err.Error(), "only supports") {
		t.Errorf("wg with no args should pass validation, got: %v", err)
	}
}

func TestRunWg_InvalidSubcommand(t *testing.T) {
	// wg with other subcommands should be rejected
	_, err := RunWg("set", "wg0", "peer", "xxxxx")
	if err == nil {
		t.Error("Expected error for wg with unsupported subcommand, got nil")
	}
	if !strings.Contains(err.Error(), "only supports 'show'") {
		t.Errorf("Expected error about unsupported wg command, got: %v", err)
	}
}

// Test wg-quick command validation
func TestRunWgQuick_Valid(t *testing.T) {
	// Create a temporary config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")
	err := os.WriteFile(configPath, []byte("[Interface]\n"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// This may fail if wg-quick is not installed, but should pass validation
	_, err = RunWgQuick("up", configPath)
	// We accept either success or execution failure (if wg-quick not installed)
	// but validation errors should NOT occur for valid inputs
	if err != nil && (strings.Contains(err.Error(), "must be 'up' or 'down'") ||
		strings.Contains(err.Error(), "must be absolute") ||
		strings.Contains(err.Error(), ".conf extension")) {
		t.Errorf("wg-quick up should pass validation, got: %v", err)
	}
}

func TestRunWgQuick_Down_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")
	err := os.WriteFile(configPath, []byte("[Interface]\n"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err = RunWgQuick("down", configPath)
	// We accept either success or execution failure (if wg-quick not installed)
	// but validation errors should NOT occur for valid inputs
	if err != nil && (strings.Contains(err.Error(), "must be 'up' or 'down'") ||
		strings.Contains(err.Error(), "must be absolute") ||
		strings.Contains(err.Error(), ".conf extension")) {
		t.Errorf("wg-quick down should pass validation, got: %v", err)
	}
}

func TestRunWgQuick_RelativePath(t *testing.T) {
	_, err := RunWgQuick("up", "relative/path/test.conf")
	if err == nil {
		t.Error("Expected error for relative path, got nil")
	}
	if !strings.Contains(err.Error(), "must be absolute") {
		t.Errorf("Expected error about absolute path requirement, got: %v", err)
	}
}

func TestRunWgQuick_NonConfFile(t *testing.T) {
	_, err := RunWgQuick("up", "/tmp/test.txt")
	if err == nil {
		t.Error("Expected error for non-.conf file, got nil")
	}
	if !strings.Contains(err.Error(), ".conf extension") {
		t.Errorf("Expected error about .conf extension, got: %v", err)
	}
}

func TestRunWgQuick_PathTraversal(t *testing.T) {
	_, err := RunWgQuick("up", "/tmp/../etc/test.conf")
	if err == nil {
		t.Error("Expected error for path traversal attempt, got nil")
	}
	if !strings.Contains(err.Error(), "invalid config file path") {
		t.Errorf("Expected error about invalid path, got: %v", err)
	}
}

// Test sudo requirement
func TestRunWgQuick_SudoNotAvailable(t *testing.T) {
	// Temporarily modify PATH to hide sudo
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", oldPath)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.conf")

	_, err := RunWgQuick("up", configPath)
	if err == nil {
		t.Error("Expected error when sudo is not available, got nil")
	}
	if !strings.Contains(err.Error(), "sudo is required but not available") {
		t.Errorf("Expected error about sudo not available, got: %v", err)
	}
}

func TestRunWg_SudoNotAvailable(t *testing.T) {
	// Temporarily modify PATH to hide sudo
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", oldPath)

	_, err := RunWg("show")
	if err == nil {
		t.Error("Expected error when sudo is not available, got nil")
	}
	if !strings.Contains(err.Error(), "sudo is required but not available") {
		t.Errorf("Expected error about sudo not available, got: %v", err)
	}
}
