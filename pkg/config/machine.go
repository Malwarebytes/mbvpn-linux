package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type MachineIdProvider interface {
	Get() (string, error)
}

type EtcFileMachineIdProvider struct{}

func NewEtcFileMachineIdProvider() MachineIdProvider {
	return &EtcFileMachineIdProvider{}
}

func (cp *EtcFileMachineIdProvider) Get() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}

	id := string(data)
	id = strings.TrimSuffix(id, "\n")

	return id, nil
}

type ConfigFileMachineIdProvider struct{}

func NewConfigFileMachineIdProvider() MachineIdProvider {
	return &ConfigFileMachineIdProvider{}
}

func (cp *ConfigFileMachineIdProvider) Get() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".config", "mbvpn")
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	filePath := filepath.Join(configPath, "machine-id")
	
	newUUID := ""
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0{ // if file does not exists or empty, create it
		newUUID = uuid.New().String()
		if err := os.WriteFile(filePath, []byte(newUUID), 0600); err != nil {
			return "", fmt.Errorf("failed to write machine-id file: %w", err)
		}
	} else {
		newUUID = string(data)
	}
	
	return strings.TrimSuffix(newUUID, "\n"), nil
}