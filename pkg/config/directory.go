package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DirectoryProvider defines the interface for getting application directories and file paths.
// This abstraction allows for easy testing by injecting test directories instead of using
// the actual user config directory.
type DirectoryProvider interface {
	// GetConfigDir returns the main config directory (e.g., ~/.config/mbvpn)
	GetConfigDir() (string, error)

	// GetServersDir returns the servers directory (e.g., ~/.config/mbvpn/servers)
	GetServersDir() (string, error)

	// GetConfigFile returns the config file path (e.g., ~/.config/mbvpn/config.yml)
	GetConfigFile() (string, error)

	// GetMachineIDFile returns the machine ID file path (e.g., ~/.config/mbvpn/machine-id)
	GetMachineIDFile() (string, error)

	// GetServersFile returns the servers JSON file path (e.g., ~/.config/mbvpn/servers.json)
	GetServersFile() (string, error)
}

// directoryProvider is the implementation that uses a base directory
type directoryProvider struct {
	baseDir string
}

// NewDefaultDirectoryProvider creates a provider using the OS user config directory.
// Returns an error if the user config directory cannot be determined.
func NewDefaultDirectoryProvider() (DirectoryProvider, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}
	return NewDirectoryProvider(userConfigDir), nil
}

// NewDirectoryProvider creates a provider with a custom base directory.
// This is primarily used for testing with temporary directories.
func NewDirectoryProvider(baseDir string) DirectoryProvider {
	return &directoryProvider{baseDir: baseDir}
}

// GetConfigDir returns the main config directory
func (d *directoryProvider) GetConfigDir() (string, error) {
	return filepath.Join(d.baseDir, "mbvpn"), nil
}

// GetServersDir returns the servers directory
func (d *directoryProvider) GetServersDir() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "servers"), nil
}

// GetConfigFile returns the config file path
func (d *directoryProvider) GetConfigFile() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yml"), nil
}

// GetMachineIDFile returns the machine ID file path
func (d *directoryProvider) GetMachineIDFile() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "machine-id"), nil
}

// GetServersFile returns the servers JSON file path
func (d *directoryProvider) GetServersFile() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "servers.json"), nil
}
