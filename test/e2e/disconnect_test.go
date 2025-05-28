//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestDisconnectCommandBasic(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command without any specific server
	output, err := execCommand(t, []string{"disconnect"}, homeDir)
	
	// Disconnect command may succeed or fail depending on system state
	// but it should handle the case gracefully
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Disconnect command crashed unexpectedly: %s", output)
		}
		
		// Expected errors include sudo requirements, wg command issues, or no active connections
		validErrors := []string{"sudo", "wg-quick", "permission", "no active", "not found"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Disconnect command failed with unexpected error: %s", output)
		}
	}
	
	// Output should not be empty in most cases
	if strings.TrimSpace(output) == "" && err != nil {
		t.Errorf("Disconnect command produced empty output with error: %v", err)
	}
}

func TestDisconnectCommandWithServer(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command with a specific server
	output, err := execCommand(t, []string{"disconnect", "us-newyork-1"}, homeDir)
	
	// Disconnect command may succeed or fail depending on system state
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Disconnect command crashed unexpectedly: %s", output)
		}
		
		// Expected errors include sudo requirements, wg command issues, or server not found
		validErrors := []string{"sudo", "wg-quick", "permission", "not found", "no active"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Disconnect command failed with unexpected error: %s", output)
		}
	}
}

func TestDisconnectCommandWithHelp(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command with help flag
	output, err := execCommand(t, []string{"disconnect", "--help"}, homeDir)
	
	if err != nil {
		t.Errorf("Disconnect command with help flag failed with error: %v", err)
	}
	
	// Should show help text
	if !strings.Contains(output, "Stops the VPN connection") {
		t.Errorf("Expected help text in output, got: %s", output)
	}
	
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
}

func TestDisconnectCommandShortAlias(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command using short alias
	output, err := execCommand(t, []string{"d", "--help"}, homeDir)
	
	if err != nil {
		t.Errorf("Disconnect command short alias failed with error: %v", err)
	}
	
	// Should show help text for disconnect command
	if !strings.Contains(output, "Stops the VPN connection") {
		t.Errorf("Expected disconnect command help text, got: %s", output)
	}
}

func TestDisconnectCommandOutputStructure(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command
	output, err := execCommand(t, []string{"disconnect"}, homeDir)
	
	// Command may fail due to network requirements or system state
	if err != nil {
		// Output should be properly formatted (no error traces unless debug mode)
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for i, line := range lines {
			// Lines should not contain Go stack traces or panic information
			if strings.Contains(line, "panic:") || strings.Contains(line, "goroutine") {
				t.Errorf("Disconnect output contains unexpected error trace on line %d: %s", i+1, line)
			}
		}
	}
}

func TestDisconnectCommandFormattingConsistency(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command multiple times to ensure consistent output
	output1, err1 := execCommand(t, []string{"disconnect"}, homeDir)
	output2, err2 := execCommand(t, []string{"disconnect"}, homeDir)
	
	// Both should have the same error status
	if (err1 == nil) != (err2 == nil) {
		t.Errorf("Disconnect command should have consistent error status across runs")
	}
	
	// Both should produce similar output structure (allowing for timestamps)
	// Strip timestamps to compare structure
	cleanOutput1 := strings.ReplaceAll(output1, "time=", "TIME=")
	cleanOutput2 := strings.ReplaceAll(output2, "time=", "TIME=")
	
	// Remove specific timestamp values but keep the general structure
	lines1 := strings.Split(strings.TrimSpace(cleanOutput1), "\n")
	lines2 := strings.Split(strings.TrimSpace(cleanOutput2), "\n")
	
	if len(lines1) != len(lines2) {
		t.Errorf("Disconnect command should produce consistent output structure across runs.\nFirst lines: %d\nSecond lines: %d", len(lines1), len(lines2))
	}
}

func TestDisconnectCommandWithInvalidServer(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command with invalid server name
	output, err := execCommand(t, []string{"disconnect", "invalid-server-name"}, homeDir)
	
	// Command should handle invalid server gracefully
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Disconnect command crashed unexpectedly with invalid server: %s", output)
		}
		
		// Expected errors include server not found, sudo requirements, or wg command issues
		validErrors := []string{"not found", "sudo", "wg-quick", "permission", "no active"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Disconnect command failed with unexpected error for invalid server: %s", output)
		}
	}
}

func TestDisconnectCommandNoActiveConnection(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command when no connection is active
	output, err := execCommand(t, []string{"disconnect"}, homeDir)
	
	// Command may succeed (no-op) or fail with "no active connection"
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Disconnect command crashed unexpectedly: %s", output)
		}
		
		// Expected errors include no active connections, sudo requirements, or wg command issues
		validErrors := []string{"no active", "not found", "sudo", "wg-quick", "permission"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Disconnect command failed with unexpected error when no connection active: %s", output)
		}
	}
}

func TestDisconnectCommandWithExtraArguments(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command with multiple arguments (should use first one)
	output, err := execCommand(t, []string{"disconnect", "server1", "extra-arg"}, homeDir)
	
	// Command should handle extra arguments gracefully (usually ignores them)
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Disconnect command crashed unexpectedly with extra args: %s", output)
		}
		
		// Expected errors include server not found, sudo requirements, or wg command issues
		validErrors := []string{"not found", "sudo", "wg-quick", "permission", "no active"}
		hasValidError := false
		for _, validErr := range validErrors {
			if strings.Contains(strings.ToLower(output), validErr) {
				hasValidError = true
				break
			}
		}
		
		if !hasValidError {
			t.Errorf("Disconnect command failed with unexpected error with extra args: %s", output)
		}
	}
}

func TestDisconnectCommandWithDebugFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the disconnect command with debug flag
	output, err := execCommand(t, []string{"disconnect", "--debug"}, homeDir)
	
	// Command may fail due to system requirements
	if err != nil {
		// Check if it's a reasonable error (not a panic or crash)
		if strings.Contains(output, "panic:") || strings.Contains(output, "goroutine") {
			t.Errorf("Disconnect command crashed unexpectedly with debug flag: %s", output)
		}
	}
	
	// Debug output should contain some information
	// At minimum, it should show some output (even if it's an error)
	if strings.TrimSpace(output) == "" {
		t.Errorf("Disconnect command with debug flag produced no output")
	}
}