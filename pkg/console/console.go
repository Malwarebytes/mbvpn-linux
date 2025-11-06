package console

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Command execution timeout
const commandTimeout = 30 * time.Second

// validateWgQuickArgs validates arguments for wg-quick command
func validateWgQuickArgs(action string, configPath string) error {
	// Validate action
	if action != "up" && action != "down" {
		return fmt.Errorf("wg-quick action must be 'up' or 'down', got '%s'", action)
	}

	// Validate config file path - must be absolute and .conf extension
	if !filepath.IsAbs(configPath) {
		return fmt.Errorf("config file path must be absolute")
	}
	if filepath.Ext(configPath) != ".conf" {
		return fmt.Errorf("config file must have .conf extension")
	}

	// Ensure no path traversal attempts
	cleanPath := filepath.Clean(configPath)
	if cleanPath != configPath {
		return fmt.Errorf("invalid config file path")
	}

	return nil
}

// validateWgArgs validates arguments for wg command
func validateWgArgs(args []string) error {
	// wg show has no arguments or just "show"
	if len(args) > 1 || (len(args) == 1 && args[0] != "show") {
		return fmt.Errorf("wg only supports 'show' command")
	}
	return nil
}

// RunWgQuick executes wg-quick command with sudo.
// Only allows 'up' and 'down' actions with validated config file paths.
func RunWgQuick(action string, configPath string) (string, error) {
	// Validate arguments
	if err := validateWgQuickArgs(action, configPath); err != nil {
		log.Warnf("Argument validation failed for wg-quick: %v", err)
		return "", err
	}

	// Check if sudo is available
	_, err := exec.LookPath("sudo")
	if err != nil {
		return "", fmt.Errorf("sudo is required but not available: %w", err)
	}

	var stdout strings.Builder
	var stderr strings.Builder

	// Create context with timeout to prevent hanging commands
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	// Log privileged command execution for audit purposes
	log.Infof("Executing privileged command: wg-quick %s %s", action, configPath)
	cmd := exec.CommandContext(ctx, "sudo", "wg-quick", action, configPath)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// Log detailed error for debugging but return sanitized error to caller
		log.Debugf("Command 'wg-quick %s %s' failed - stderr: %v, stdout: %v", action, configPath, stderr.String(), stdout.String())

		// Check if context deadline exceeded (timeout)
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command execution timed out after %v", commandTimeout)
		}

		return stdout.String(), fmt.Errorf("command execution failed: %w", err)
	}

	return stdout.String(), nil
}

// RunWg executes wg command with sudo.
// Only allows 'show' subcommand or no arguments.
func RunWg(args ...string) (string, error) {
	// Validate arguments
	if err := validateWgArgs(args); err != nil {
		log.Warnf("Argument validation failed for wg: %v", err)
		return "", err
	}

	// Check if sudo is available
	_, err := exec.LookPath("sudo")
	if err != nil {
		return "", fmt.Errorf("sudo is required but not available: %w", err)
	}

	var stdout strings.Builder
	var stderr strings.Builder

	// Create context with timeout to prevent hanging commands
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	// Build command arguments
	cmdArgs := []string{"wg"}
	if len(args) > 0 {
		cmdArgs = append(cmdArgs, args...)
	}

	// Log privileged command execution for audit purposes
	log.Infof("Executing privileged command: wg %v", args)
	cmd := exec.CommandContext(ctx, "sudo", cmdArgs...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// Log detailed error for debugging but return sanitized error to caller
		log.Debugf("Command 'wg %v' failed - stderr: %v, stdout: %v", args, stderr.String(), stdout.String())

		// Check if context deadline exceeded (timeout)
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command execution timed out after %v", commandTimeout)
		}

		return stdout.String(), fmt.Errorf("command execution failed: %w", err)
	}

	return stdout.String(), nil
}