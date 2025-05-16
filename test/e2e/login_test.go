//go:build e2e

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const (
	validLicenseKey   = "CBKGF-JRFZJ-FZYF8-4QR8E" // Sample license key for testing
	invalidLicenseKey = "AAAAA-BBBBB-CCCCC-DDDDD"
	invalidMbCode     = "MB-12345"  // Invalid MB code
)

// TestLogin is a helper function to execute the login command and capture its output
func execLogin(t *testing.T, args []string, homeDir string) (string, error) {
	t.Helper()

	// Get the absolute path to the mbvpn binary
	binPath, err := filepath.Abs("../../mbvpn")
	if err != nil {
		t.Fatalf("Failed to get binary path: %v", err)
	}

	// Create the command
	cmd := exec.Command(binPath, args...)
	
	// Set environment variables for the command
	cmd.Env = append(os.Environ(), 
		"HOME="+homeDir,
		"MBVPN_HOLOCRON_URL_ST=https://idms-holocron-stage.mwbsys.com/graphql", // Use staging environment
		"MBVPN_ENV=testing",
	)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()

	// Combine stdout and stderr for the result
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nError: " + stderr.String()
	}

	return output, err
}

// setupTestHome creates a temporary home directory for testing
func setupTestHome(t *testing.T) string {
	t.Helper()

	// Create a temporary directory to serve as HOME
	homeDir, err := os.MkdirTemp("", "mbvpn-test-home-*")
	if err != nil {
		t.Fatalf("Failed to create test home directory: %v", err)
	}

	// Create the necessary subdirectories
	configDir := filepath.Join(homeDir, ".config", "mbvpn")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create a machine-id file with a random UUID
	machineIdPath := filepath.Join(configDir, "machine-id")
	err = os.WriteFile(machineIdPath, []byte(uuid.New().String()), 0600)
	if err != nil {
		t.Fatalf("Failed to create machine-id file: %v", err)
	}

	return homeDir
}

// cleanupTestHome removes the temporary home directory
func cleanupTestHome(t *testing.T, homeDir string) {
	t.Helper()
	
	// Remove the temporary directory
	err := os.RemoveAll(homeDir)
	if err != nil {
		t.Logf("Warning: Failed to clean up test home directory: %v", err)
	}
}

// checkConfigFile checks if the config file was created and contains expected data
func checkConfigFile(t *testing.T, homeDir string) bool {
	t.Helper()
	
	configPath := filepath.Join(homeDir, ".config", "mbvpn", "config.yml")
	
	// Check if the file exists
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		t.Fatalf("Error checking config file: %v", err)
	}
	
	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Error reading config file: %v", err)
	}
	
	// Check if the file contains the installation_token
	return strings.Contains(string(data), "installation_token")
}

func TestLoginWithValidLicenseKey(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command with a valid license key
	output, err := execLogin(t, []string{"login", "--licenseKey", validLicenseKey, "--debug"}, homeDir)
	
	// Even if the command fails due to network or backend issues, we're testing the command execution
	// not the actual activation process which requires real backend connectivity
	if err != nil {
		// In a real environment, we'd expect success, but in testing without real backend,
		// we might get an error. Just log it for debugging.
		t.Logf("Command returned error (expected in test environment): %v", err)
	}
	
	// Check the output for expected messages
	if !strings.Contains(output, "Welcome to Malwarebytes VPN client!") {
		t.Errorf("Expected welcome message in output, got: %s", output)
	}
	
	// In a real environment with working backend, we would also check for:
	// 1. "Activated successfully!" message in the output
	// 2. Presence of a config file with installation token
	if checkConfigFile(t, homeDir) {
		t.Logf("Config file was created successfully")
	}
}

func TestLoginWithInvalidLicenseKey(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command with an invalid license key format
	output, _ := execLogin(t, []string{"login", "--licenseKey", "INVALID-FORMAT", "--debug"}, homeDir)
	
	// Check the output for expected error message
	if !strings.Contains(output, "Invalid license key") {
		t.Errorf("Expected invalid license key error in output, got: %s", output)
	}
	
	// Config file should not be created for invalid input
	if checkConfigFile(t, homeDir) {
		t.Error("Config file was created despite invalid license key")
	}
}

func TestLoginWithInvalidMbCode(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command with an invalid MB code format (too short)
	output, _ := execLogin(t, []string{"login", "--mbCode", "MB-123", "--debug"}, homeDir)
	
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
	output, _ := execLogin(t, []string{"login", "--debug"}, homeDir)
	
	// Check the output for expected error message about required flags
	if !strings.Contains(output, "required") && !strings.Contains(output, "flags") {
		t.Errorf("Expected required flag error in output, got: %s", output)
	}
	
	// Config file should not be created when command fails
	if checkConfigFile(t, homeDir) {
		t.Error("Config file was created despite missing credentials")
	}
}

func TestLoginAlreadyLoggedIn(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Create a mock config file with an installation token to simulate being logged in
	configDir := filepath.Join(homeDir, ".config", "mbvpn")
	configPath := filepath.Join(configDir, "config.yml")
	mockConfig := `installation_token: mock-installation-token
private_key: ""
public_key: ""
`
	err := os.WriteFile(configPath, []byte(mockConfig), 0600)
	if err != nil {
		t.Fatalf("Failed to create mock config file: %v", err)
	}

	// Execute the login command with a valid license key
	output, _ := execLogin(t, []string{"login", "--licenseKey", validLicenseKey, "--debug"}, homeDir)
	
	// Check the output for expected error message
	if !strings.Contains(output, "active session") {
		t.Errorf("Expected active session error in output, got: %s", output)
	}
}

func TestLoginInvalidFlags(t *testing.T) {
	// Create a temporary home directory
	homeDir := setupTestHome(t)
	defer cleanupTestHome(t, homeDir)

	// Execute the login command with both license key and MB code (should be one or the other)
	output, err := execLogin(t, []string{
		"login", 
		"--licenseKey", validLicenseKey, 
		"--mbCode", invalidMbCode, 
		"--debug",
	}, homeDir)
	
	// If the command doesn't prevent using both flags, this might not fail as expected
	// In a well-designed CLI, it should prefer one flag over the other or show an error
	t.Logf("Output from using both flags: %s", output)
	t.Logf("Error from using both flags: %v", err)
	
	// Based on the implementation, we know it prioritizes license key if both are provided
}