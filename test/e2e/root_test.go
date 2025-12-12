//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestRootCommandNoArgs(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command without arguments
	output, err := execCommand(t, []string{}, homeDir)
	
	if err != nil {
		t.Errorf("Root command without arguments failed with error: %v", err)
	}
	
	// Should show help text when no arguments are provided
	if !strings.Contains(output, "MBVPN is a command-line VPN client") {
		t.Errorf("Expected help text in output, got: %s", output)
	}
	
	// Should show usage information
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
	
	// Should show available commands
	if !strings.Contains(output, "Available Commands:") {
		t.Errorf("Expected available commands section in output, got: %s", output)
	}
}

func TestRootCommandWithHelp(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command with help flag
	output, err := execCommand(t, []string{"--help"}, homeDir)
	
	if err != nil {
		t.Errorf("Root command with help flag failed with error: %v", err)
	}
	
	// Should show detailed help text
	if !strings.Contains(output, "MBVPN is a command-line VPN client for Malwarebytes VPN service") {
		t.Errorf("Expected short description in help output, got: %s", output)
	}
	
	// Should show long description
	if !strings.Contains(output, "MBVPN is a command-line VPN client for Malwarebytes VPN service") {
		t.Errorf("Expected long description in help output, got: %s", output)
	}
	
	// Should mention license requirement
	if !strings.Contains(output, "Requires a valid Malwarebytes license") {
		t.Errorf("Expected license requirement in help output, got: %s", output)
	}
	
	// Should show global flags
	if !strings.Contains(output, "--debug") {
		t.Errorf("Expected debug flag in help output, got: %s", output)
	}
	
	if !strings.Contains(output, "--trace") {
		t.Errorf("Expected trace flag in help output, got: %s", output)
	}
}

func TestRootCommandWithDebugFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command with debug flag
	output, err := execCommand(t, []string{"--debug"}, homeDir)
	
	if err != nil {
		t.Errorf("Root command with debug flag failed with error: %v", err)
	}
	
	// Debug mode indicator only appears when a subcommand is executed
	// Root command with no subcommand shows help, so no debug indicator
	// Should still show help text since no command was specified
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
	
	// Should show the debug flag in help
	if !strings.Contains(output, "--debug") {
		t.Errorf("Expected debug flag in help output, got: %s", output)
	}
}

func TestRootCommandWithTraceFlag(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command with trace flag
	output, err := execCommand(t, []string{"--trace"}, homeDir)
	
	if err != nil {
		t.Errorf("Root command with trace flag failed with error: %v", err)
	}
	
	// Trace mode indicator only appears when a subcommand is executed
	// Root command with no subcommand shows help, so no trace indicator
	// Should still show help text since no command was specified
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage information in output, got: %s", output)
	}
	
	// Should show the trace flag in help
	if !strings.Contains(output, "--trace") {
		t.Errorf("Expected trace flag in help output, got: %s", output)
	}
}

func TestRootCommandAvailableCommands(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command
	output, err := execCommand(t, []string{}, homeDir)
	
	if err != nil {
		t.Errorf("Root command failed with error: %v", err)
	}
	
	// Check that all expected commands are listed
	expectedCommands := []string{
		"cities",
		"connect",
		"countries",
		"disconnect",
		"login",
		"logout",
		"servers",
		"status",
		"version",
	}
	
	for _, cmd := range expectedCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Expected command '%s' to be listed in available commands, got: %s", cmd, output)
		}
	}
}

func TestRootCommandInvalidSubcommand(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command with invalid subcommand
	output, err := execCommand(t, []string{"invalid-command"}, homeDir)
	
	// Should fail with error
	if err == nil {
		t.Errorf("Root command with invalid subcommand should fail, but succeeded with output: %s", output)
	}
	
	// Should show error message about unknown command
	if !strings.Contains(output, "unknown command") && !strings.Contains(output, "Error: unknown command") {
		t.Errorf("Expected error message about unknown command, got: %s", output)
	}
}

func TestRootCommandGlobalFlagsWithSubcommand(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute a subcommand with global debug flag
	output, err := execCommand(t, []string{"--debug", "version"}, homeDir)
	
	if err != nil {
		t.Errorf("Root command with debug flag and version subcommand failed with error: %v", err)
	}
	
	// Should show debug mode indicator
	if !strings.Contains(output, "=== Running in debug mode ===") {
		t.Errorf("Expected debug mode indicator in output, got: %s", output)
	}
	
	// Should show version information (format: vX.Y.Z)
	if !strings.Contains(output, "v0.0.") {
		t.Errorf("Expected version information in output, got: %s", output)
	}
}

func TestRootCommandFlagOrder(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Test flag before subcommand
	output1, err1 := execCommand(t, []string{"--debug", "version"}, homeDir)
	
	// Test flag after subcommand
	output2, err2 := execCommand(t, []string{"version", "--debug"}, homeDir)
	
	// Both should work (Cobra supports both flag positions)
	if err1 != nil {
		t.Errorf("Flag before subcommand failed: %v", err1)
	}
	
	if err2 != nil {
		t.Errorf("Flag after subcommand failed: %v", err2)
	}
	
	// Both should show debug mode
	if !strings.Contains(output1, "=== Running in debug mode ===") {
		t.Errorf("Flag before subcommand should enable debug mode, got: %s", output1)
	}
	
	if !strings.Contains(output2, "=== Running in debug mode ===") {
		t.Errorf("Flag after subcommand should enable debug mode, got: %s", output2)
	}
}

func TestRootCommandOutputStructure(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the root command
	output, err := execCommand(t, []string{}, homeDir)
	
	if err != nil {
		t.Errorf("Root command failed with error: %v", err)
	}
	
	// Output should not be empty
	if strings.TrimSpace(output) == "" {
		t.Errorf("Root command produced empty output")
	}
	
	// Output should be properly formatted (no error traces)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for i, line := range lines {
		// Lines should not contain Go stack traces or panic information
		if strings.Contains(line, "panic:") || strings.Contains(line, "goroutine") {
			t.Errorf("Root output contains unexpected error trace on line %d: %s", i+1, line)
		}
	}
	
	// Should have multiple lines (help text is multi-line)
	if len(lines) < 5 {
		t.Errorf("Expected multi-line help output, got %d lines: %s", len(lines), output)
	}
}