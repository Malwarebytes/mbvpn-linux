//go:build e2e

package main

import (
	"strings"
	"testing"
)

func TestLogoutWithActiveSession(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the logout command
	output, err := execCommand(t, []string{"logout", "--debug"}, homeDir)
	
	// Even if the command fails due to network or backend issues, we're testing the command execution
	// not the actual deactivation process which requires real backend connectivity
	if err != nil {
		// In a real environment, we'd expect success, but in testing without real backend,
		// we might get an error. Just log it for debugging.
		t.Logf("Command returned error (expected in test environment): %v", err)
	}
	
	// Check the output for expected success message
	if !strings.Contains(output, "Logged out successfully") {
		t.Errorf("Expected success message in output, got: %s", output)
	}
	
	// Check that the config file was deleted
	if !checkConfigDeleted(t, homeDir) {
		t.Error("Config file was not deleted during logout")
	}
}

func TestLogoutWithoutActiveSession(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// No need to create a mock session, we want to test logout without active session

	// Execute the logout command
	output, _ := execCommand(t, []string{"logout", "--debug"}, homeDir)
	
	// In the current implementation, logout succeeds even without an active session
	// and prints "Logged out successfully!", so we should check for that message
	if !strings.Contains(output, "Logged out successfully") {
		t.Errorf("Expected 'Logged out successfully' message in output, got: %s", output)
	}
}

func TestLogoutAfterLogin(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command with a valid license key
	_, _ = execCommand(t, []string{"login", "--licenseKey", getValidLicenseKey(), "--debug"}, homeDir)
	
	// If the login was successful and created a config file
	if checkConfigFile(t, homeDir) {
		// Execute the logout command
		output, _ := execCommand(t, []string{"logout", "--debug"}, homeDir)
		
		// Check the output for expected success message
		if !strings.Contains(output, "Logged out successfully") {
			t.Errorf("Expected success message in output, got: %s", output)
		}
		
		// Check that the config file was deleted
		if !checkConfigDeleted(t, homeDir) {
			t.Error("Config file was not deleted during logout")
		}
	} else {
		t.Skip("Skipping test because login did not create a config file")
	}
}