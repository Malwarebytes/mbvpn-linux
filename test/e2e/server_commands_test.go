//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// testServerCommandWithSession tests a server-related command with an active session
func testServerCommandWithSession(t *testing.T, command string, expectedContent []string, unexpectedContent []string) {
	t.Helper()
	
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session to simulate active login
	createMockSession(t, homeDir)

	// Execute the command
	output, err := execCommand(t, []string{command}, homeDir)
	
	if err != nil {
		t.Errorf("%s command with active session failed with error: %v", command, err)
	}

	// Check for expected content
	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected '%s' in %s output, got: %s", expected, command, output)
		}
	}

	// Check for unexpected content
	for _, unexpected := range unexpectedContent {
		if strings.Contains(output, unexpected) {
			t.Errorf("Unexpected '%s' in %s output, got: %s", unexpected, command, output)
		}
	}
}

// testServerCommandNoSession tests a server-related command without an active session
func testServerCommandNoSession(t *testing.T, command string) {
	t.Helper()
	
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the command without creating a session
	output, err := execCommand(t, []string{command}, homeDir)
	
	// Command exits with error code 1 but prints error message
	if err == nil {
		t.Errorf("%s command without session should fail, but succeeded with output: %s", command, output)
	}

	// Should show unauthorized error message
	if !strings.Contains(output, "no active session") {
		t.Errorf("Expected 'no active session' error in %s output, got: %s", command, output)
	}

	// Should suggest login command
	if !strings.Contains(output, "login") {
		t.Errorf("Expected login suggestion in %s output, got: %s", command, output)
	}
}

// testServerCommandWithDebugFlag tests a server-related command with debug flag
func testServerCommandWithDebugFlag(t *testing.T, command string) {
	t.Helper()
	
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the command with debug flag
	output, err := execCommand(t, []string{command, "--debug"}, homeDir)
	
	if err != nil {
		t.Errorf("%s command with debug flag failed with error: %v", command, err)
	}

	// Should show debug mode indicator
	if !strings.Contains(output, "=== Running in debug mode ===") {
		t.Errorf("Expected debug mode indicator in %s output, got: %s", command, output)
	}
}

// testServerCommandHelp tests help output for a server-related command
func testServerCommandHelp(t *testing.T, command string, expectedDescription string) {
	t.Helper()
	
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the command with help flag
	output, err := execCommand(t, []string{command, "--help"}, homeDir)
	
	if err != nil {
		t.Errorf("%s command help failed with error: %v", command, err)
	}

	// Should show usage information
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in %s help output, got: %s", command, output)
	}

	// Should show command description (check if part of description is present)
	if !strings.Contains(output, expectedDescription) {
		t.Errorf("Expected description containing '%s' in %s help output, got: %s", expectedDescription, command, output)
	}

	// Should mention session requirement
	if !strings.Contains(output, "active session") {
		t.Errorf("Expected session requirement mention in %s help output, got: %s", command, output)
	}
}

// testServerCommandOutputStructure tests the basic output structure
func testServerCommandOutputStructure(t *testing.T, command string) {
	t.Helper()
	
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the command
	output, err := execCommand(t, []string{command}, homeDir)
	
	if err != nil {
		t.Errorf("%s command failed with error: %v", command, err)
	}

	// Output should not be empty
	if strings.TrimSpace(output) == "" {
		t.Errorf("%s command produced empty output", command)
	}

	// Output should be properly formatted (no error traces)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for i, line := range lines {
		// Lines should not contain Go stack traces or panic information
		if strings.Contains(line, "panic:") || strings.Contains(line, "goroutine") {
			t.Errorf("%s output contains unexpected error trace on line %d: %s", command, i+1, line)
		}
	}
}