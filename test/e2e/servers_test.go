//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestServersCommandWithSession(t *testing.T) {
	expectedContent := []string{
		// Since servers shows the most comprehensive view (cities + server details)
		// We expect it to potentially show server names, cities, and connection info
	}
	
	unexpectedContent := []string{
		"no active session",
		"Error:",
		"panic:",
	}
	
	testServerCommandWithSession(t, "servers", expectedContent, unexpectedContent)
}

func TestServersCommandNoSession(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the servers command without creating a session
	output, err := execCommand(t, []string{"servers"}, homeDir)
	
	// Servers command works without session (uses cached data)
	if err != nil {
		t.Errorf("Servers command failed unexpectedly: %v", err)
	}

	// Should show server information even without session
	if strings.TrimSpace(output) == "" {
		t.Errorf("Servers command should show server data even without session")
	}
}

func TestServersCommandWithDebugFlag(t *testing.T) {
	testServerCommandWithDebugFlag(t, "servers")
}

func TestServersCommandHelp(t *testing.T) {
	testServerCommandHelp(t, "servers", "Displays a comprehensive list")
}

func TestServersCommandOutputStructure(t *testing.T) {
	testServerCommandOutputStructure(t, "servers")
}

func TestServersCommandComprehensiveOutput(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the servers command
	output, err := execCommand(t, []string{"servers"}, homeDir)
	
	if err != nil {
		t.Errorf("Servers command failed with error: %v", err)
	}

	// Servers command should show the most detailed information
	// It calls vpn.Servers(true, true) which means cities=true, details=true
	// This should be the most comprehensive output compared to countries and cities
	
	// The output should not be empty for a properly functioning command
	if strings.TrimSpace(output) == "" {
		t.Errorf("Servers command produced empty output")
	}

	// Should not contain error messages if session is active
	if strings.Contains(output, "no active session") {
		t.Errorf("Servers command with active session should not show session error")
	}
}

func TestServersCommandWithTraceFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the servers command with trace flag
	output, err := execCommand(t, []string{"servers", "--trace"}, homeDir)
	
	if err != nil {
		t.Errorf("Servers command with trace flag failed with error: %v", err)
	}

	// Should show trace mode indicator
	if !strings.Contains(output, "=== Running in trace mode ===") {
		t.Errorf("Expected trace mode indicator in servers output, got: %s", output)
	}
}

func TestServersCommandInvalidFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the servers command with invalid flag
	output, err := execCommand(t, []string{"servers", "--invalid-flag"}, homeDir)
	
	// Should fail with error
	if err == nil {
		t.Errorf("Servers command with invalid flag should fail, but succeeded with output: %s", output)
	}

	// Should show error message about unknown flag
	if !strings.Contains(output, "unknown flag") && !strings.Contains(output, "Error: unknown flag") {
		t.Errorf("Expected error message about unknown flag, got: %s", output)
	}
}