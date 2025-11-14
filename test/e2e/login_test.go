//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestLoginWithInvalidMbCode(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command with an invalid MB code format (too short)
	output, _ := execCommand(t, []string{"login", "--code", "MB-123", "--debug"}, homeDir)

	// Check the output for expected error message
	if !strings.Contains(output, "Invalid MB-code") {
		t.Errorf("Expected invalid MB code error in output, got: %s", output)
	}

	// Config file should not be created for invalid input
	if checkConfigFile(t, homeDir) {
		t.Error("Config file was created despite invalid MB code")
	}
}

func TestLoginNoCredentials(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command without providing credentials
	output, _ := execCommand(t, []string{"login", "--debug"}, homeDir)

	// Check the output for expected error message about required flags
	if !strings.Contains(output, "required") && !strings.Contains(output, "flags") {
		t.Errorf("Expected required flag error in output, got: %s", output)
	}

	// Config file should not be created when command fails
	if checkConfigFile(t, homeDir) {
		t.Error("Config file was created despite missing credentials")
	}
}
