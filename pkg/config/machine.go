package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type MachineIdProvider interface {
	Get() (string, error)
}

type ConfigFileMachineIdProvider struct{
	dirProvider DirectoryProvider
}

func NewConfigFileMachineIdProvider(dirProvider DirectoryProvider) MachineIdProvider {
	return &ConfigFileMachineIdProvider{
		dirProvider: dirProvider,
	}
}

func (cp *ConfigFileMachineIdProvider) Get() (string, error) {
	filePath, err := cp.dirProvider.GetMachineIDFile()
	if err != nil {
		return "", fmt.Errorf("failed to get machine ID file path: %w", err)
	}

	// Ensure directory exists
	configDir := filepath.Dir(filePath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	newUUID := ""
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 { // if file does not exists or empty, create it
		newUUID = uuid.New().String()
		if err := os.WriteFile(filePath, []byte(newUUID), 0600); err != nil {
			return "", fmt.Errorf("failed to write machine-id file: %w", err)
		}
	} else {
		newUUID = string(data)
	}

	return strings.TrimSuffix(newUUID, "\n"), nil
}
