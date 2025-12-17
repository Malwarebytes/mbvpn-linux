//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestStatusCommandNotLoggedIn(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the status command without being logged in
	output, err := execCommand(t, []string{"status"}, homeDir)
	
	// Status command works without session and shows network status
	if err != nil {
		t.Errorf("Status command failed unexpectedly: %v", err)
	}
	
	// Should show network status information
	if strings.TrimSpace(output) == "" {
		t.Errorf("Status command should show network status, got empty output")
	}
	
	// Should show VPN status (enabled/disabled)
	if !strings.Contains(output, "VPN enabled:") {
		t.Errorf("Expected VPN status information, got: %s", output)
	}
}

func TestStatusCommandNotConnected(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session to simulate being logged in but not connected
	createMockSession(t, homeDir)

	// Execute the status command
	output, err := execCommand(t, []string{"status"}, homeDir)
	
	// Status command may fail due to sudo requirements for wg commands
	if err != nil {
		// If it fails, it should be due to sudo requirements, not session issues
		if !strings.Contains(output, "sudo") {
			t.Errorf("Status command failed for unexpected reason: %v, output: %s", err, output)
		}
		return
	}
	
	// If it succeeds, check that output indicates not connected status
	if !strings.Contains(output, "VPN enabled: false") && !strings.Contains(output, "No active connections") {
		t.Errorf("Expected disconnected status message, got: %s", output)
	}
}

func TestStatusCommandWithDebugFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the status command with debug flag
	output, err := execCommand(t, []string{"status", "--debug"}, homeDir)
	
	// Status command may fail due to sudo requirements
	if err != nil {
		// If it fails, it should be due to sudo requirements
		if !strings.Contains(output, "sudo") {
			t.Errorf("Status command with debug flag failed for unexpected reason: %v", err)
		}
	}
	
	// Debug output should contain additional information
	// At minimum, it should show some output (even if it's an error)
	if strings.TrimSpace(output) == "" {
		t.Errorf("Status command with debug flag produced no output")
	}
}

func TestStatusCommandWithHelp(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the status command with help flag
	output, err := execCommand(t, []string{"status", "--help"}, homeDir)
	
	if err != nil {
		t.Errorf("Status command with help flag failed with error: %v", err)
	}
	
	// Should show help text (check for key parts that exist in the actual output)
	if !strings.Contains(output, "current VPN connection status") {
		t.Errorf("Expected help text in output, got: %s", output)
	}
	
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
	
	// Should mention detailed information
	if !strings.Contains(output, "Displays detailed information about the current VPN connection status") {
		t.Errorf("Expected detailed description in help output, got: %s", output)
	}
}

func TestStatusCommandFormattingConsistency(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the status command multiple times to ensure consistent output
	output1, err1 := execCommand(t, []string{"status"}, homeDir)
	output2, err2 := execCommand(t, []string{"status"}, homeDir)
	
	// Both should have the same error status
	if (err1 == nil) != (err2 == nil) {
		t.Errorf("Status command should have consistent error status across runs")
	}
	
	// Both should produce similar output structure (allowing for timestamps)
	// Strip timestamps to compare structure
	cleanOutput1 := strings.ReplaceAll(output1, "time=", "TIME=")
	cleanOutput2 := strings.ReplaceAll(output2, "time=", "TIME=")
	
	// Remove specific timestamp values but keep the general structure
	lines1 := strings.Split(strings.TrimSpace(cleanOutput1), "\n")
	lines2 := strings.Split(strings.TrimSpace(cleanOutput2), "\n")
	
	if len(lines1) != len(lines2) {
		t.Errorf("Status command should produce consistent output structure across runs.\nFirst lines: %d\nSecond lines: %d", len(lines1), len(lines2))
	}
}

func TestStatusCommandInvalidArguments(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the status command with invalid arguments
	output, err := execCommand(t, []string{"status", "invalid-arg"}, homeDir)
	
	// Command should handle extra arguments gracefully
	// Most CLI tools ignore extra arguments for status commands
	if err != nil {
		// If it fails, it should be due to sudo requirements or session issues
		if !strings.Contains(output, "not logged in") &&
		   !strings.Contains(output, "No active session") &&
		   !strings.Contains(output, "sudo") {
			t.Errorf("Status command failed for unexpected reason with invalid args: %v, output: %s", err, output)
		}
	}
}

func TestStatusCommandOutputStructure(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute the status command
	output, err := execCommand(t, []string{"status"}, homeDir)
	
	// Status command may fail due to sudo requirements
	if err != nil {
		// If it fails, it should be due to sudo requirements
		if !strings.Contains(output, "sudo") {
			t.Errorf("Status command failed for unexpected reason: %v", err)
		}
	}
	
	// Output should not be empty
	if strings.TrimSpace(output) == "" {
		t.Errorf("Status command produced empty output")
	}
	
	// Output should be properly formatted (no error traces unless debug mode)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for i, line := range lines {
		// Lines should not contain Go stack traces or panic information
		if strings.Contains(line, "panic:") || strings.Contains(line, "goroutine") {
			t.Errorf("Status output contains unexpected error trace on line %d: %s", i+1, line)
		}
	}
}