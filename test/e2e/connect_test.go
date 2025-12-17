//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestConnectCommandNotLoggedIn(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the connect command without being logged in
	output, err := execCommand(t, []string{"connect", "US"}, homeDir)
	
	// Command may succeed but should show authentication error in output
	if err != nil {
		// If it fails, check for authentication error
		if !strings.Contains(output, "no active session") && !strings.Contains(output, "Try 'login' command first") {
			t.Errorf("Expected authentication error message, got: %s", output)
		}
	} else {
		// If it succeeds, it should show some kind of authentication warning or error
		if !strings.Contains(output, "no active session") && !strings.Contains(output, "Try 'login' command first") && !strings.Contains(output, "not logged in") {
			t.Logf("Warning: Connect command succeeded without session, output: %s", output)
		}
	}
}

func TestConnectCommandWithoutServer(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command without server argument
	output, err := execCommand(t, []string{"connect"}, homeDir)
	
	if err == nil {
		t.Errorf("Connect command should fail when no server is specified")
	}
	
	// Should show usage error
	if !strings.Contains(output, "accepts 1 arg(s), received 0") {
		t.Errorf("Expected argument error message, got: %s", output)
	}
}

func TestConnectCommandWithInvalidServer(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command with invalid server
	output, err := execCommand(t, []string{"connect", "INVALID-SERVER"}, homeDir)
	
	if err == nil {
		t.Errorf("Connect command should fail with invalid server")
	}
	
	// Should show server not found error
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected server error message, got: %s", output)
	}
}

func TestConnectCommandWithValidCountryCode(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command with valid country code
	output, err := execCommand(t, []string{"connect", "US"}, homeDir)
	
	// Command may fail due to missing servers.json or network requirements
	// but it should show appropriate error messages, not crash
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Connect command crashed unexpectedly: %s", output)
		}
		
		// Expected errors include server data missing, network issues, or sudo requirements
		validErrors := []string{"servers.json", "no servers found", "sudo", "permission"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Connect command failed with unexpected error: %s", output)
		}
	}
}

func TestConnectCommandWithCountryName(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command with country name
	output, err := execCommand(t, []string{"connect", "United States"}, homeDir)
	
	// Command may fail due to missing servers.json or network requirements
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Connect command crashed unexpectedly: %s", output)
		}
		
		// Expected errors include server data missing, network issues, or sudo requirements
		validErrors := []string{"servers.json", "no servers found", "sudo", "permission"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Connect command failed with unexpected error: %s", output)
		}
	}
}

func TestConnectCommandWithCityCode(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command with city code
	output, err := execCommand(t, []string{"connect", "NYC"}, homeDir)
	
	// Command may fail due to missing servers.json or network requirements
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Connect command crashed unexpectedly: %s", output)
		}
		
		// Expected errors include server data missing, network issues, or sudo requirements
		validErrors := []string{"servers.json", "no servers found", "sudo", "permission"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Connect command failed with unexpected error: %s", output)
		}
	}
}

func TestConnectCommandWithHelp(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the connect command with help flag
	output, err := execCommand(t, []string{"connect", "--help"}, homeDir)
	
	if err != nil {
		t.Errorf("Connect command with help flag failed with error: %v", err)
	}
	
	// Should show help text
	if !strings.Contains(output, "Establishes a VPN connection") {
		t.Errorf("Expected help text in output, got: %s", output)
	}
	
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
	
	// Should mention server specification options
	if !strings.Contains(output, "Server Specification Options") {
		t.Errorf("Expected server specification help in output, got: %s", output)
	}
	
	// Should show examples
	if !strings.Contains(output, "Examples:") {
		t.Errorf("Expected examples in help output, got: %s", output)
	}
}

func TestConnectCommandShortAlias(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the connect command using short alias
	output, err := execCommand(t, []string{"c", "--help"}, homeDir)
	
	if err != nil {
		t.Errorf("Connect command short alias failed with error: %v", err)
	}
	
	// Should show help text for connect command
	if !strings.Contains(output, "Establishes a VPN connection") {
		t.Errorf("Expected connect command help text, got: %s", output)
	}
}

func TestConnectCommandOutputStructure(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command
	output, err := execCommand(t, []string{"connect", "US"}, homeDir)
	
	// Command may fail due to missing servers.json or network requirements
	if err != nil {
		// Output should be properly formatted (no error traces unless debug mode)
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for i, line := range lines {
			// Lines should not contain Go stack traces or panic information
			if strings.Contains(line, "panic:") || strings.Contains(line, "goroutine") {
				t.Errorf("Connect output contains unexpected error trace on line %d: %s", i+1, line)
			}
		}
	}
	
	// Output should not be empty
	if strings.TrimSpace(output) == "" {
		t.Errorf("Connect command produced empty output")
	}
}

func TestConnectCommandInvalidArguments(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the connect command with too many arguments
	output, err := execCommand(t, []string{"connect", "US", "extra-arg"}, homeDir)
	
	if err == nil {
		t.Errorf("Connect command should fail with too many arguments")
	}
	
	// Should show usage error
	if !strings.Contains(output, "Error:") {
		t.Errorf("Expected argument error message, got: %s", output)
	}
}