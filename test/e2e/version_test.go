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

	// Check that output contains version information
	if !strings.Contains(output, "Version:") {
		t.Errorf("Expected 'Version:' in output, got: %s", output)
	}

	// Check that output contains environment information
	if !strings.Contains(output, "Environment:") {
		t.Errorf("Expected 'Environment:' in output, got: %s", output)
	}

	// Check that output contains build type information
	if !strings.Contains(output, "Build Type:") {
		t.Errorf("Expected 'Build Type:' in output, got: %s", output)
	}

	// Check that output contains debug mode information
	if !strings.Contains(output, "Debug Mode:") {
		t.Errorf("Expected 'Debug Mode:' in output, got: %s", output)
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

	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have exactly 4 lines of output
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines of output, got %d lines: %s", len(lines), output)
	}

	// Check version format (should be like "Version: 0.0.1+0")
	versionLine := lines[0]
	if !strings.HasPrefix(versionLine, "Version: ") {
		t.Errorf("First line should start with 'Version: ', got: %s", versionLine)
	}

	versionValue := strings.TrimPrefix(versionLine, "Version: ")
	if !strings.Contains(versionValue, ".") || !strings.Contains(versionValue, "+") {
		t.Errorf("Version should be in format 'major.minor.patch+build', got: %s", versionValue)
	}

	// Check environment line
	envLine := lines[1]
	if !strings.HasPrefix(envLine, "Environment: ") {
		t.Errorf("Second line should start with 'Environment: ', got: %s", envLine)
	}

	// Check debug mode line
	debugLine := lines[3]
	if !strings.HasPrefix(debugLine, "Debug Mode: ") {
		t.Errorf("Fourth line should start with 'Debug Mode: ', got: %s", debugLine)
	}

	debugValue := strings.TrimPrefix(debugLine, "Debug Mode: ")
	if debugValue != "true" && debugValue != "false" {
		t.Errorf("Debug Mode should be 'true' or 'false', got: %s", debugValue)
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

	// Should still contain all version information
	expectedFields := []string{"Version:", "Environment:", "Build Type:", "Debug Mode:"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("Expected '%s' in output, got: %s", field, output)
		}
	}

	// When debug flag is used, debug mode should be true
	if !strings.Contains(output, "Debug Mode: true") {
		t.Errorf("Expected 'Debug Mode: true' when using --debug flag, got: %s", output)
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
	if !strings.Contains(output, "Display the application version, build environment and build type") {
		t.Errorf("Expected help text in output, got: %s", output)
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
}
