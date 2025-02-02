package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var Debug bool

type Config struct {
	InstallationToken string `yaml:"installation_token"`
}

func StoreInstallationToken(token string) {
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
	f, err := os.Create(configPath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.Encode(cfg)
}

func GetInstallationToken() (string, error) {
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

func DeleteConfig() error {
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

func GetMachineId() (string, error) {
  data, err := os.ReadFile("/etc/machine-id")
  if err != nil {
    return "", err
  }

  id := string(data)
  id = strings.TrimSuffix(id, "\n")

  return id, nil
}
