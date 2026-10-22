//go:build e2e

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const (
	invalidMbCode = "MB-12345"
)

// Accept development builds and SemVer releases, including prerelease/build metadata.
var versionPattern = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(\.(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)

func validVersion(version string) bool {
	return version == "vdev" || versionPattern.MatchString(version)
}

func assertVersionOutput(t *testing.T, output string) {
	t.Helper()
	if !validVersion(strings.TrimSpace(output)) {
		t.Errorf("Expected vdev or a semantic version, got: %s", output)
	}
}

func testConfigDir(homeDir string) string {
	return filepath.Join(homeDir, ".config", "mbvpn")
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
	configDir := testConfigDir(homeDir)
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

	configPath := filepath.Join(testConfigDir(homeDir), "config.yml")

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

// checkConfigDeleted checks if the config file was properly deleted after logout
func checkConfigDeleted(t *testing.T, homeDir string) bool {
	t.Helper()

	configPath := filepath.Join(testConfigDir(homeDir), "config.yml")

	// Check if the file exists
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}

// createMockSession creates a mock config file to simulate an active session
func createMockSession(t *testing.T, homeDir string) {
	t.Helper()

	configDir := testConfigDir(homeDir)
	configPath := filepath.Join(configDir, "config.yml")
	mockConfig := `installation_token: mock-installation-token
private_key: ""
public_key: ""
`
	err := os.WriteFile(configPath, []byte(mockConfig), 0600)
	if err != nil {
		t.Fatalf("Failed to create mock config file: %v", err)
	}
}

// execCommand is a helper function to execute any mbvpn command and capture its output
func execCommand(t *testing.T, args []string, homeDir string) (string, error) {
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
		"XDG_CONFIG_HOME="+filepath.Join(homeDir, ".config"),
		"MBVPN_HOLOCRON_URL=https://idms-holocron-stage.mwbsys.com/graphql", // Use staging environment
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
