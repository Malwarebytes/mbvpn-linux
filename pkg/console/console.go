package console

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RunCmd executes a terminal command, optionally with sudo.
func RunCmd(sudo bool, command string, args ...string) (string, error) {
	var cmd *exec.Cmd

	var stdout strings.Builder
	var stderr strings.Builder

	if sudo {
		// Check if sudo is available
		_, err := exec.LookPath("sudo")
		if err != nil {
			log.Debugln("Sudo not available, running without sudo...")
			cmd = exec.Command(command, args...)
		} else {
			cmd = exec.Command("sudo", append([]string{command}, args...)...)
		}
	} else {
		cmd = exec.Command(command, args...)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("command execution failed: %w, stderr: %v, stdout: %v", err, stderr.String(), stdout.String())
	}

	return stdout.String(), nil
}

func WgShow() (string, error) {
	var stdout strings.Builder
	var stderr strings.Builder

	cmd := exec.Command("sudo", "wg", "show")

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("wg show command execution failed: %w, stderr: %v, stdout: %v", err, stderr.String(), stdout.String())
	}

	return stdout.String(), nil
}

func WgUp(cfgPath string) error {
	cleanedPath, err := sanitizeWgConfigPath(cfgPath)
	if err != nil {
		return fmt.Errorf("invalid WireGuard config path: %w", err)
	}

	var stdout strings.Builder
	var stderr strings.Builder

	cmd := exec.Command("sudo", "wg-quick", "up", cleanedPath)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("wg-quick up command execution failed: %w, stderr: %v, stdout: %v", err, stderr.String(), stdout.String())
	}

	return nil
}

func WgDown(cfgPath string) error {
	cleanedPath, err := sanitizeWgConfigPath(cfgPath)
	if err != nil {
		return fmt.Errorf("invalid WireGuard config path: %w", err)
	}

	var stdout strings.Builder
	var stderr strings.Builder

	cmd := exec.Command("sudo", "wg-quick", "down", cleanedPath)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("wg-quick down command execution failed: %w, stderr: %v, stdout: %v", err, stderr.String(), stdout.String())
	}

	return nil
}

func sanitizeWgConfigPath(cfgPath string) (string, error) {
	// Clean the path (removes .., redundant separators, etc.)
	cleanPath := filepath.Clean(cfgPath)

	// Expand home directory if present
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	} else if strings.HasPrefix(cleanPath, "~/") {
		cleanPath = filepath.Join(home, cleanPath[2:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Verify the file exists and is a regular file
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("config file not accessible: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("config path is not a regular file")
	}

	// Ensure it's within the expected config directory
	expectedBase := filepath.Join(home, ".config", "mbvpn", "servers")
	if !strings.HasPrefix(absPath, expectedBase) {
		return "", fmt.Errorf("config file must be within %s", expectedBase)
	}

	// Verify file extension (WireGuard configs must be .conf)
	if filepath.Ext(absPath) != ".conf" {
		return "", fmt.Errorf("config file must have .conf extension")
	}

	return absPath, nil
}
