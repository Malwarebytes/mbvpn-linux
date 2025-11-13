package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InstallationToken string `yaml:"installation_token"`
	PrivateKey string `yaml:"private_key"`
	PublicKey string `yaml:"public_key"`
}

type ConfigProvider interface {
	StoreInstallationToken(token string) error
	GetInstallationToken() (string, error)
	DeleteConfig() error
	Get() (Config, error)
	StoreData(publicKey string, privateKey string) error
	update(cfg Config) error
}

type YamlConfigProvider struct{
	dirProvider DirectoryProvider
}

func NewYamlConfigProvider(dirProvider DirectoryProvider) ConfigProvider {
	return &YamlConfigProvider{
		dirProvider: dirProvider,
	}
}

func (cp *YamlConfigProvider) StoreInstallationToken(token string) error {
	cfg, _ := cp.Get()
	cfg.InstallationToken = token
	err := cp.update(cfg)
	if err != nil {
		return fmt.Errorf("failed to store installation token: %w", err)
	}
	return nil
}

func (cp *YamlConfigProvider) GetInstallationToken() (string, error) {
	cfg, err := cp.Get()
	if err != nil {
		return "", fmt.Errorf("failed to get installation token: %w", err)
	}

	return cfg.InstallationToken, nil
}

func (cp *YamlConfigProvider) DeleteConfig() error {
	configPath, err := cp.dirProvider.GetConfigFile()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	err = os.Remove(configPath)
	if err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	return nil
}

func (cp *YamlConfigProvider) update(cfg Config) error {
	configPath, err := cp.dirProvider.GetConfigFile()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file for writing: %w", err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config data: %w", err)
	}

	return nil
}

func (cp *YamlConfigProvider) StoreData(publicKey string, privateKey string) error {
	cfg, err := cp.Get()
	if err != nil {
		return fmt.Errorf("failed to get existing config: %w", err)
	}

	cfg.PrivateKey = privateKey
	cfg.PublicKey = publicKey 

	err = cp.update(cfg)
	if err != nil {
		return fmt.Errorf("failed to update config with new data: %w", err)
	}

	return nil
}


func (cp *YamlConfigProvider) Get() (Config, error) {
	configPath, err := cp.dirProvider.GetConfigFile()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get config file path: %w", err)
	}

	// Try to open the file
	f, err := os.Open(configPath)

	// If file doesn't exist, create it with default config
	if os.IsNotExist(err) {
		// Ensure directory exists
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return Config{}, fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create default config
		defaultConfig := Config{
			// Set your default values here
			// For example:
			// ServerURL: "default.server.com",
			// Port: 51820,
		}

		// Create and write to the file
		f, err := os.Create(configPath)
		if err != nil {
			return Config{}, fmt.Errorf("failed to create config file: %w", err)
		}
		defer f.Close()

		encoder := yaml.NewEncoder(f)
		if err := encoder.Encode(defaultConfig); err != nil {
			return Config{}, fmt.Errorf("failed to write default config: %w", err)
		}

		return defaultConfig, nil
	} else if err != nil {
		// Handle other errors
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	// If we got here, the file exists and is open, so decode it
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
}

