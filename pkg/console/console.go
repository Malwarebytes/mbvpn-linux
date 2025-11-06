package console

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// allowedCommands defines the commands that can be executed
var allowedCommands = map[string]bool{
	"wg":       true,
	"wg-quick": true,
}

// allowedSudoCommands defines commands that can run with sudo
var allowedSudoCommands = map[string]bool{
	"wg":       true,
	"wg-quick": true,
}

// Command execution timeout
const commandTimeout = 30 * time.Second

// validateCommand checks if a command is allowed to be executed
func validateCommand(sudo bool, command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Extract just the command name (in case it's a path)
	cmdName := filepath.Base(command)

	// Check if command is in the allowed list
	if !allowedCommands[cmdName] {
		return fmt.Errorf("command '%s' is not allowed", cmdName)
	}

	// Check if command is allowed with sudo
	if sudo && !allowedSudoCommands[cmdName] {
		return fmt.Errorf("command '%s' is not allowed to run with sudo", cmdName)
	}

	return nil
}

// validateArgs performs basic validation on command arguments
func validateArgs(command string, args []string) error {
	cmdName := filepath.Base(command)

	switch cmdName {
	case "wg-quick":
		// wg-quick expects: up/down <config-file>
		if len(args) < 2 {
			return fmt.Errorf("wg-quick requires action (up/down) and config file")
		}
		action := args[0]
		if action != "up" && action != "down" {
			return fmt.Errorf("wg-quick action must be 'up' or 'down', got '%s'", action)
		}
		// Validate config file path - must be absolute and .conf extension
		configPath := args[1]
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

	case "wg":
		// wg show has no arguments or just "show"
		if len(args) > 1 || (len(args) == 1 && args[0] != "show") {
			return fmt.Errorf("wg only supports 'show' command")
		}

	default:
		// For any other allowed commands, ensure no shell metacharacters
		invalidChars := regexp.MustCompile(`[;&|$\x60<>(){}[\]!*?~]`)
		for _, arg := range args {
			if invalidChars.MatchString(arg) {
				return fmt.Errorf("argument contains invalid characters: %s", arg)
			}
		}
	}

	return nil
}

// RunCmd executes a terminal command, optionally with sudo.
// Only allows pre-approved commands (wg, wg-quick) with validated arguments.
func RunCmd(sudo bool, command string, args ...string) (string, error) {
	// Validate command is allowed
	if err := validateCommand(sudo, command); err != nil {
		log.Warnf("Command validation failed: %v", err)
		return "", err
	}

	// Validate arguments
	if err := validateArgs(command, args); err != nil {
		log.Warnf("Argument validation failed for command '%s': %v", command, err)
		return "", err
	}

	var cmd *exec.Cmd
	var stdout strings.Builder
	var stderr strings.Builder

	// Create context with timeout to prevent hanging commands
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	if sudo {
		// Check if sudo is available
		_, err := exec.LookPath("sudo")
		if err != nil {
			return "", fmt.Errorf("sudo is required but not available: %w", err)
		}

		// Log privileged command execution for audit purposes
		log.Infof("Executing privileged command: %s %v", command, args)
		cmd = exec.CommandContext(ctx, "sudo", append([]string{command}, args...)...)
	} else {
		log.Debugf("Executing command: %s %v", command, args)
		cmd = exec.CommandContext(ctx, command, args...)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Log detailed error for debugging but return sanitized error to caller
		log.Debugf("Command '%s %v' failed - stderr: %v, stdout: %v", command, args, stderr.String(), stdout.String())

		// Check if context deadline exceeded (timeout)
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command execution timed out after %v", commandTimeout)
		}

		return stdout.String(), fmt.Errorf("command execution failed: %w", err)
	}

	return stdout.String(), nil
}