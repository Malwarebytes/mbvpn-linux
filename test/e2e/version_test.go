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

	assertVersionOutput(t, output)
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

	assertVersionOutput(t, output)
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

	// Validate the version separately from the debug banner.
	assertVersionOutput(t, strings.TrimPrefix(output, "=== Running in debug mode ===\n"))

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
