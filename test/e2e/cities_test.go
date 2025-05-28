//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestCitiesCommandWithSession(t *testing.T) {
	expectedContent := []string{
		// Cities command shows city information but not full server details
		// vpn.Servers(true, false) means cities=true, details=false
	}
	
	unexpectedContent := []string{
		"no active session",
		"Error:",
		"panic:",
	}
	
	testServerCommandWithSession(t, "cities", expectedContent, unexpectedContent)
}

func TestCitiesCommandNoSession(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the cities command without creating a session
	output, err := execCommand(t, []string{"cities"}, homeDir)
	
	// Cities command works without session (uses cached data)
	if err != nil {
		t.Errorf("Cities command failed unexpectedly: %v", err)
	}

	// Should show city information even without session
	if strings.TrimSpace(output) == "" {
		t.Errorf("Cities command should show city data even without session")
	}
}

func TestCitiesCommandWithDebugFlag(t *testing.T) {
	testServerCommandWithDebugFlag(t, "cities")
}

func TestCitiesCommandHelp(t *testing.T) {
	testServerCommandHelp(t, "cities", "Displays a list of cities")
}

func TestCitiesCommandOutputStructure(t *testing.T) {
	testServerCommandOutputStructure(t, "cities")
}

func TestCitiesCommandSpecificOutput(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the cities command
	output, err := execCommand(t, []string{"cities"}, homeDir)
	
	if err != nil {
		t.Errorf("Cities command failed with error: %v", err)
	}

	// Cities command should show city information without full server details
	// It calls vpn.Servers(true, false) which means cities=true, details=false
	// This should show more detail than countries but less than servers
	
	// The output should not be empty for a properly functioning command
	if strings.TrimSpace(output) == "" {
		t.Errorf("Cities command produced empty output")
	}

	// Should not contain error messages if session is active
	if strings.Contains(output, "no active session") {
		t.Errorf("Cities command with active session should not show session error")
	}
}

func TestCitiesCommandWithTraceFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the cities command with trace flag
	output, err := execCommand(t, []string{"cities", "--trace"}, homeDir)
	
	if err != nil {
		t.Errorf("Cities command with trace flag failed with error: %v", err)
	}

	// Should show trace mode indicator
	if !strings.Contains(output, "=== Running in trace mode ===") {
		t.Errorf("Expected trace mode indicator in cities output, got: %s", output)
	}
}

func TestCitiesCommandInvalidFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the cities command with invalid flag
	output, err := execCommand(t, []string{"cities", "--invalid-flag"}, homeDir)
	
	// Should fail with error
	if err == nil {
		t.Errorf("Cities command with invalid flag should fail, but succeeded with output: %s", output)
	}

	// Should show error message about unknown flag
	if !strings.Contains(output, "unknown flag") && !strings.Contains(output, "Error: unknown flag") {
		t.Errorf("Expected error message about unknown flag, got: %s", output)
	}
}

func TestCitiesCommandFlagOrder(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Test flag before command
	output1, err1 := execCommand(t, []string{"--debug", "cities"}, homeDir)
	
	// Test flag after command
	output2, err2 := execCommand(t, []string{"cities", "--debug"}, homeDir)
	
	// Both should work (Cobra supports both flag positions)
	if err1 != nil {
		t.Errorf("Flag before cities command failed: %v", err1)
	}
	
	if err2 != nil {
		t.Errorf("Flag after cities command failed: %v", err2)
	}
	
	// Both should show debug mode
	if !strings.Contains(output1, "=== Running in debug mode ===") {
		t.Errorf("Flag before cities command should enable debug mode, got: %s", output1)
	}
	
	if !strings.Contains(output2, "=== Running in debug mode ===") {
		t.Errorf("Flag after cities command should enable debug mode, got: %s", output2)
	}
}

func TestCitiesCommandArguments(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the cities command with unexpected arguments
	output, err := execCommand(t, []string{"cities", "extra-arg"}, homeDir)
	
	// Cities command doesn't expect arguments, but Cobra might ignore them
	// The command should still execute successfully
	if err != nil {
		t.Errorf("Cities command with extra arguments failed: %v", err)
	}

	// Should not contain session errors
	if strings.Contains(output, "no active session") {
		t.Errorf("Cities command with active session should not show session error")
	}
}