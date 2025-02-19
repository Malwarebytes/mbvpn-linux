package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func Debug() bool {
	return testing.Testing() && testing.Verbose()
}

type Config struct {
	InstallationToken string `yaml:"installation_token"`
}

type ConfigProvider interface {
	StoreInstallationToken(token string)
	GetInstallationToken() (string, error)
	DeleteConfig() error
}

type YamlConfigProvider struct{}

func NewYamlConfigProvider() ConfigProvider {
	return &YamlConfigProvider{}
}

type MachineIdProvider interface {
	Get() (string, error)
}

type EtcFileMachineIdProvider struct{}

func NewEtcFileMachineIdProvider() MachineIdProvider {
	return &EtcFileMachineIdProvider{}
}

func (cp *YamlConfigProvider) StoreInstallationToken(token string) {
	cfg := &Config{
		InstallationToken: token,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	configDir := filepath.Join(home, ".config", "mbvpn")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		log.Panic(err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.Encode(cfg)
}

func (cp *YamlConfigProvider) GetInstallationToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".config", "mbvpn", "config.yaml")
	f, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
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

	configPath := filepath.Join(home, ".config", "mbvpn", "config.yaml")

	err = os.Remove(configPath)
	if err != nil {
		return err
	}

	return nil
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
