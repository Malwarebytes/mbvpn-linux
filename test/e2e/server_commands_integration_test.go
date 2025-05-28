//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestServerCommandsIntegration(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute all three commands
	serversOutput, serversErr := execCommand(t, []string{"servers"}, homeDir)
	countriesOutput, countriesErr := execCommand(t, []string{"countries"}, homeDir)
	citiesOutput, citiesErr := execCommand(t, []string{"cities"}, homeDir)

	// All commands should execute without fatal errors when session is active
	if serversErr != nil {
		t.Errorf("Servers command failed: %v", serversErr)
	}
	if countriesErr != nil {
		t.Errorf("Countries command failed: %v", countriesErr)
	}
	if citiesErr != nil {
		t.Errorf("Cities command failed: %v", citiesErr)
	}

	// All outputs should be non-empty
	if strings.TrimSpace(serversOutput) == "" {
		t.Error("Servers command produced empty output")
	}
	if strings.TrimSpace(countriesOutput) == "" {
		t.Error("Countries command produced empty output")
	}
	if strings.TrimSpace(citiesOutput) == "" {
		t.Error("Cities command produced empty output")
	}

	// None should show session errors
	commands := map[string]string{
		"servers":   serversOutput,
		"countries": countriesOutput,
		"cities":    citiesOutput,
	}

	for cmd, output := range commands {
		if strings.Contains(output, "no active session") {
			t.Errorf("%s command with active session should not show session error", cmd)
		}
	}
}

func TestServerCommandsConsistentBehaviorWithoutSession(t *testing.T) {
	// Create a temporary home directory without session
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute all three commands without session
	serversOutput, serversErr := execCommand(t, []string{"servers"}, homeDir)
	countriesOutput, countriesErr := execCommand(t, []string{"countries"}, homeDir)
	citiesOutput, citiesErr := execCommand(t, []string{"cities"}, homeDir)

	// All commands should work without session (using cached data)
	if serversErr != nil {
		t.Errorf("Servers command failed: %v", serversErr)
	}
	if countriesErr != nil {
		t.Errorf("Countries command failed: %v", countriesErr)
	}
	if citiesErr != nil {
		t.Errorf("Cities command failed: %v", citiesErr)
	}

	// All should show server data even without session
	commands := map[string]string{
		"servers":   serversOutput,
		"countries": countriesOutput,
		"cities":    citiesOutput,
	}

	for cmd, output := range commands {
		if strings.TrimSpace(output) == "" {
			t.Errorf("%s command should show data even without session", cmd)
		}
	}
}

func TestServerCommandsWithDebugConsistency(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock session
	createMockSession(t, homeDir)

	// Execute all three commands with debug flag
	serversOutput, serversErr := execCommand(t, []string{"servers", "--debug"}, homeDir)
	countriesOutput, countriesErr := execCommand(t, []string{"countries", "--debug"}, homeDir)
	citiesOutput, citiesErr := execCommand(t, []string{"cities", "--debug"}, homeDir)

	// All commands should handle debug flag consistently
	if serversErr != nil {
		t.Errorf("Servers command with debug failed: %v", serversErr)
	}
	if countriesErr != nil {
		t.Errorf("Countries command with debug failed: %v", countriesErr)
	}
	if citiesErr != nil {
		t.Errorf("Cities command with debug failed: %v", citiesErr)
	}

	// All should show debug mode indicator
	commands := map[string]string{
		"servers":   serversOutput,
		"countries": countriesOutput,
		"cities":    citiesOutput,
	}

	for cmd, output := range commands {
		if !strings.Contains(output, "=== Running in debug mode ===") {
			t.Errorf("%s command should show debug mode indicator, got: %s", cmd, output)
		}
	}
}

func TestServerCommandsHelpConsistency(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute all three commands with help flag
	serversOutput, serversErr := execCommand(t, []string{"servers", "--help"}, homeDir)
	countriesOutput, countriesErr := execCommand(t, []string{"countries", "--help"}, homeDir)
	citiesOutput, citiesErr := execCommand(t, []string{"cities", "--help"}, homeDir)

	// Help should work for all commands
	if serversErr != nil {
		t.Errorf("Servers help failed: %v", serversErr)
	}
	if countriesErr != nil {
		t.Errorf("Countries help failed: %v", countriesErr)
	}
	if citiesErr != nil {
		t.Errorf("Cities help failed: %v", citiesErr)
	}

	// All help outputs should have consistent structure
	commands := map[string]string{
		"servers":   serversOutput,
		"countries": countriesOutput,
		"cities":    citiesOutput,
	}

	for cmd, output := range commands {
		if !strings.Contains(output, "Usage:") {
			t.Errorf("%s help should show usage information, got: %s", cmd, output)
		}
		if !strings.Contains(output, "active session") {
			t.Errorf("%s help should mention session requirement, got: %s", cmd, output)
		}
		if !strings.Contains(output, "Global Flags:") {
			t.Errorf("%s help should show global flags, got: %s", cmd, output)
		}
	}
}

func TestServerCommandsInvalidFlagConsistency(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute all three commands with invalid flag
	serversOutput, serversErr := execCommand(t, []string{"servers", "--invalid"}, homeDir)
	countriesOutput, countriesErr := execCommand(t, []string{"countries", "--invalid"}, homeDir)
	citiesOutput, citiesErr := execCommand(t, []string{"cities", "--invalid"}, homeDir)

	// All commands should fail with invalid flag
	if serversErr == nil {
		t.Error("Servers command should fail with invalid flag")
	}
	if countriesErr == nil {
		t.Error("Countries command should fail with invalid flag")
	}
	if citiesErr == nil {
		t.Error("Cities command should fail with invalid flag")
	}

	// All should show consistent error messages
	commands := map[string]string{
		"servers":   serversOutput,
		"countries": countriesOutput,
		"cities":    citiesOutput,
	}

	for cmd, output := range commands {
		if !strings.Contains(output, "unknown flag") && !strings.Contains(output, "Error: unknown flag") {
			t.Errorf("%s command should show unknown flag error, got: %s", cmd, output)
		}
	}
}