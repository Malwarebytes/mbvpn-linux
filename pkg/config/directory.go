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

// DefaultDirectoryProvider is the production implementation that uses os.UserConfigDir()
type DefaultDirectoryProvider struct{}

// NewDefaultDirectoryProvider creates a new DefaultDirectoryProvider
func NewDefaultDirectoryProvider() DirectoryProvider {
	return &DefaultDirectoryProvider{}
}

// GetConfigDir returns the main config directory
func (d *DefaultDirectoryProvider) GetConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(userConfigDir, "mbvpn"), nil
}

// GetServersDir returns the servers directory
func (d *DefaultDirectoryProvider) GetServersDir() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "servers"), nil
}

// GetConfigFile returns the config file path
func (d *DefaultDirectoryProvider) GetConfigFile() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yml"), nil
}

// GetMachineIDFile returns the machine ID file path
func (d *DefaultDirectoryProvider) GetMachineIDFile() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "machine-id"), nil
}

// GetServersFile returns the servers JSON file path
func (d *DefaultDirectoryProvider) GetServersFile() (string, error) {
	configDir, err := d.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "servers.json"), nil
}

// TestDirectoryProvider is a test implementation that uses a custom base directory
type TestDirectoryProvider struct {
	baseDir string
}

// NewTestDirectoryProvider creates a new TestDirectoryProvider with the given base directory
func NewTestDirectoryProvider(baseDir string) DirectoryProvider {
	return &TestDirectoryProvider{baseDir: baseDir}
}

// GetConfigDir returns the main config directory
func (t *TestDirectoryProvider) GetConfigDir() (string, error) {
	return filepath.Join(t.baseDir, "mbvpn"), nil
}

// GetServersDir returns the servers directory
func (t *TestDirectoryProvider) GetServersDir() (string, error) {
	configDir, err := t.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "servers"), nil
}

// GetConfigFile returns the config file path
func (t *TestDirectoryProvider) GetConfigFile() (string, error) {
	configDir, err := t.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yml"), nil
}

// GetMachineIDFile returns the machine ID file path
func (t *TestDirectoryProvider) GetMachineIDFile() (string, error) {
	configDir, err := t.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "machine-id"), nil
}

// GetServersFile returns the servers JSON file path
func (t *TestDirectoryProvider) GetServersFile() (string, error) {
	configDir, err := t.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "servers.json"), nil
}
