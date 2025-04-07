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
	IpV4 string `yaml:"ipv4"`
	IpV6 string `yaml:"ipv6"`
}

type ConfigProvider interface {
	StoreInstallationToken(token string) error
	GetInstallationToken() (string, error)
	DeleteConfig() error
	Get() (Config, error)
	StoreData(key string, ipv4 string, ipv6 string) error
	update(cfg Config) error
}

type YamlConfigProvider struct{}

func NewYamlConfigProvider() ConfigProvider {
	return &YamlConfigProvider{}
}

func (cp *YamlConfigProvider) StoreInstallationToken(token string) error {
	cfg, _ := cp.Get()
	cfg.InstallationToken = token
	err := cp.update(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cp *YamlConfigProvider) GetInstallationToken() (string, error) {
	cfg, err := cp.Get()
	if err != nil {
		return "", err
	}

	return cfg.InstallationToken, nil
}

func (cp *YamlConfigProvider) DeleteConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")

	err = os.Remove(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (cp *YamlConfigProvider) update(cfg Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")
	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cp *YamlConfigProvider) StoreData(key string, ipv4 string, ipv6 string) error {
	cfg, err := cp.Get()
	if err != nil {
		return err
	}

	cfg.PrivateKey = key
	cfg.IpV4 = ipv4
	cfg.IpV6 = ipv6

	err = cp.update(cfg)
	if err != nil {
		return err
	}

	return nil
}


func (cp *YamlConfigProvider) Get() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")
	
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

