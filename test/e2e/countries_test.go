//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestCountriesCommandWithSession(t *testing.T) {
	expectedContent := []string{
		// Countries command shows only country information
		// vpn.Servers(false, false) means cities=false, details=false
	}
	
	unexpectedContent := []string{
		"no active session",
		"Error:",
		"panic:",
	}
	
	testServerCommandWithSession(t, "countries", expectedContent, unexpectedContent)
}

func TestCountriesCommandNoSession(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the countries command without creating a session
	output, err := execCommand(t, []string{"countries"}, homeDir)
	
	// Countries command works without session (uses cached data)
	if err != nil {
		t.Errorf("Countries command failed unexpectedly: %v", err)
	}

	// Should show country information even without session
	if strings.TrimSpace(output) == "" {
		t.Errorf("Countries command should show country data even without session")
	}
}

func TestCountriesCommandWithDebugFlag(t *testing.T) {
	testServerCommandWithDebugFlag(t, "countries")
}

func TestCountriesCommandHelp(t *testing.T) {
	testServerCommandHelp(t, "countries", "Displays a list of countries")
}

func TestCountriesCommandOutputStructure(t *testing.T) {
	testServerCommandOutputStructure(t, "countries")
}

func TestCountriesCommandSpecificOutput(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the countries command
	output, err := execCommand(t, []string{"countries"}, homeDir)
	
	if err != nil {
		t.Errorf("Countries command failed with error: %v", err)
	}

	// Countries command should show country-level information only
	// It calls vpn.Servers(false, false) which means no cities, no details
	// This should be the least detailed output compared to servers and cities
	
	// The output should not be empty for a properly functioning command
	if strings.TrimSpace(output) == "" {
		t.Errorf("Countries command produced empty output")
	}

	// Should not contain error messages if session is active
	if strings.Contains(output, "no active session") {
		t.Errorf("Countries command with active session should not show session error")
	}
}

func TestCountriesCommandWithTraceFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the countries command with trace flag
	output, err := execCommand(t, []string{"countries", "--trace"}, homeDir)
	
	if err != nil {
		t.Errorf("Countries command with trace flag failed with error: %v", err)
	}

	// Should show trace mode indicator
	if !strings.Contains(output, "=== Running in trace mode ===") {
		t.Errorf("Expected trace mode indicator in countries output, got: %s", output)
	}
}

func TestCountriesCommandInvalidFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the countries command with invalid flag
	output, err := execCommand(t, []string{"countries", "--invalid-flag"}, homeDir)
	
	// Should fail with error
	if err == nil {
		t.Errorf("Countries command with invalid flag should fail, but succeeded with output: %s", output)
	}

	// Should show error message about unknown flag
	if !strings.Contains(output, "unknown flag") && !strings.Contains(output, "Error: unknown flag") {
		t.Errorf("Expected error message about unknown flag, got: %s", output)
	}
}

func TestCountriesCommandFlagOrder(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Test flag before command
	output1, err1 := execCommand(t, []string{"--debug", "countries"}, homeDir)
	
	// Test flag after command
	output2, err2 := execCommand(t, []string{"countries", "--debug"}, homeDir)
	
	// Both should work (Cobra supports both flag positions)
	if err1 != nil {
		t.Errorf("Flag before countries command failed: %v", err1)
	}
	
	if err2 != nil {
		t.Errorf("Flag after countries command failed: %v", err2)
	}
	
	// Both should show debug mode
	if !strings.Contains(output1, "=== Running in debug mode ===") {
		t.Errorf("Flag before countries command should enable debug mode, got: %s", output1)
	}
	
	if !strings.Contains(output2, "=== Running in debug mode ===") {
		t.Errorf("Flag after countries command should enable debug mode, got: %s", output2)
	}
}