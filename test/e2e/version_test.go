//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the version command
	output, err := execCommand(t, []string{"version"}, homeDir)

	// Version command should always succeed
	if err != nil {
		t.Errorf("Version command failed with error: %v", err)
	}

	// Check that output starts with "v" and contains version number
	output = strings.TrimSpace(output)
	if !strings.HasPrefix(output, "v") {
		t.Errorf("Expected output to start with 'v', got: %s", output)
	}

	// Check that it contains dots (semantic version format)
	if !strings.Contains(output, ".") {
		t.Errorf("Expected version format 'vX.Y.Z', got: %s", output)
	}
}

func TestVersionCommandFormat(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the version command
	output, err := execCommand(t, []string{"version"}, homeDir)

	if err != nil {
		t.Errorf("Version command failed with error: %v", err)
	}

	output = strings.TrimSpace(output)

	// Should be single line output
	if strings.Contains(output, "\n") {
		t.Errorf("Expected single line output, got multiple lines: %s", output)
	}

	// Should start with "v"
	if !strings.HasPrefix(output, "v") {
		t.Errorf("Expected version to start with 'v', got: %s", output)
	}

	// Should be in format "vX.Y.Z"
	versionValue := strings.TrimPrefix(output, "v")
	parts := strings.Split(versionValue, ".")
	if len(parts) != 3 {
		t.Errorf("Expected version format 'vX.Y.Z', got: %s", output)
	}
}

func TestVersionCommandWithDebugFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the version command with debug flag
	output, err := execCommand(t, []string{"version", "--debug"}, homeDir)

	if err != nil {
		t.Errorf("Version command with debug flag failed with error: %v", err)
	}

	output = strings.TrimSpace(output)

	// When debug flag is used, there may be debug messages before version
	// Check if output contains version format vX.Y.Z
	if !strings.Contains(output, "v0.0.") {
		t.Errorf("Expected version format 'vX.Y.Z' in output, got: %s", output)
	}

	// Should show debug mode indicator when --debug is used
	if !strings.Contains(output, "=== Running in debug mode ===") {
		t.Errorf("Expected debug mode indicator in output, got: %s", output)
	}
}

func TestVersionCommandWithHelp(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the version command with help flag
	output, err := execCommand(t, []string{"version", "--help"}, homeDir)

	if err != nil {
		t.Errorf("Version command with help flag failed with error: %v", err)
	}

	// Should show help text
	if !strings.Contains(output, "Display the application version") {
		t.Errorf("Expected help text in output, got: %s", output)
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
}
