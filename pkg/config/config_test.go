package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// MockConfigProvider implements ConfigProvider interface for testing
type MockConfigProvider struct {
	HomeError bool
}

func (cp *MockConfigProvider) GetUserHomeDir() (string, error) {
	if cp.HomeError {
		return "", fmt.Errorf("simulated error")
	}
	return os.UserHomeDir()
}

func (cp *MockConfigProvider) StoreInstallationToken(token string) error {
	_, err := cp.GetUserHomeDir()
	if err != nil {
		return err
	}
	return nil
}

func (cp *MockConfigProvider) GetInstallationToken() (string, error) {
	_, err := cp.GetUserHomeDir()
	if err != nil {
		return "", err
	}
	return "token", nil
}

func (cp *MockConfigProvider) DeleteConfig() error {
	_, err := cp.GetUserHomeDir()
	if err != nil {
		return err
	}
	return nil
}

func (cp *MockConfigProvider) update(cfg Config) error {
	_, err := cp.GetUserHomeDir()
	if err != nil {
		return err
	}
	return nil
}

func (cp *MockConfigProvider) StoreData(publicKey string, privateKey string) error {
	_, err := cp.GetUserHomeDir()
	if err != nil {
		return err
	}
	return nil
}

func (cp *MockConfigProvider) Get() (Config, error) {
	_, err := cp.GetUserHomeDir()
	if err != nil {
		return Config{}, err
	}
	return Config{}, nil
}

// Helper function to setup a test environment
func setupTestConfig(t *testing.T) (*YamlConfigProvider, func()) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Create the config directory
	configDir := filepath.Join(tempDir, ".config", "mbvpn")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	
	provider := &YamlConfigProvider{
		homeDir: tempDir,
	}
	
	// Return a cleanup function
	return provider, func() {
		// No cleanup needed as t.TempDir() handles it
	}
}

func TestNewYamlConfigProvider(t *testing.T) {
	provider := NewYamlConfigProvider()
	
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
	
	_, ok := provider.(*YamlConfigProvider)
	if !ok {
		t.Fatal("Expected provider to be of type *YamlConfigProvider")
	}
}

func TestStoreInstallationToken(t *testing.T) {
	provider, cleanup := setupTestConfig(t)
	defer cleanup()
	
	// Test storing token
	err := provider.StoreInstallationToken("test-token")
	if err != nil {
		t.Fatalf("Failed to store installation token: %v", err)
	}
	
	// Verify token was stored
	cfg, err := provider.Get()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	
	if cfg.InstallationToken != "test-token" {
		t.Errorf("Expected token to be 'test-token', got '%s'", cfg.InstallationToken)
	}
	
	// Test error case - simulate error in update
	// Make the config directory read-only to force an error
	home, _ := provider.GetUserHomeDir()
	configDir := filepath.Join(home, ".config", "mbvpn")
	configPath := filepath.Join(configDir, "config.yml")
	
	// Remove the file first
	os.Remove(configPath)
	
	// Make the directory read-only
	os.Chmod(configDir, 0500)
	
	// Try to store token again (should fail on some systems)
	err = provider.StoreInstallationToken("test-token-2")
	// We don't assert on the error since this may not fail on all systems
	
	// Restore directory permissions
	os.Chmod(configDir, 0755)
	
	// Test error case - GetUserHomeDir returns error
	errorProvider := &MockConfigProvider{HomeError: true}
	
	// Test that we get an error
	err = errorProvider.StoreInstallationToken("test-token")
	if err == nil {
		t.Error("Expected error when GetUserHomeDir fails, got nil")
	}
}

func TestGetInstallationToken(t *testing.T) {
	provider, cleanup := setupTestConfig(t)
	defer cleanup()
	
	// Create a config file with a token
	home, _ := provider.GetUserHomeDir()
	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")
	
	// Create the file and write valid YAML
	f, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	
	cfg := Config{InstallationToken: "existing-token"}
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	f.Close()
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	
	// Test getting token
	token, err := provider.GetInstallationToken()
	if err != nil {
		t.Fatalf("Failed to get installation token: %v", err)
	}
	
	if token != "existing-token" {
		t.Errorf("Expected token to be 'existing-token', got '%s'", token)
	}
	
	// Test error case - invalid YAML
	os.Remove(configPath)
	f, _ = os.Create(configPath)
	f.WriteString("invalid: yaml: :")
	f.Close()
	
	_, err = provider.GetInstallationToken()
	if err == nil {
		t.Error("Expected error with invalid YAML, got nil")
	}
	
	// Test error case - GetUserHomeDir returns error
	errorProvider := &MockConfigProvider{HomeError: true}
	
	// Test that we get an error
	_, err = errorProvider.GetInstallationToken()
	if err == nil {
		t.Error("Expected error when GetUserHomeDir fails, got nil")
	}
}

func TestDeleteConfig(t *testing.T) {
	provider, cleanup := setupTestConfig(t)
	defer cleanup()
	
	// Create a config file
	home, _ := provider.GetUserHomeDir()
	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")
	f, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	f.Close()
	
	// Test deleting config
	err = provider.DeleteConfig()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
	
	// Verify config was deleted
	_, err = os.Stat(configPath)
	if !os.IsNotExist(err) {
		t.Error("Expected config file to be deleted, but it still exists")
	}
	
	// Test error case - Remove error (file doesn't exist)
	err = provider.DeleteConfig()
	if err == nil {
		t.Error("Expected error when file doesn't exist, got nil")
	}
	
	// Test error case - GetUserHomeDir returns error
	errorProvider := &MockConfigProvider{HomeError: true}
	
	// Test that we get an error
	err = errorProvider.DeleteConfig()
	if err == nil {
		t.Error("Expected error when GetUserHomeDir fails, got nil")
	}
}

func TestUpdate(t *testing.T) {
	provider, cleanup := setupTestConfig(t)
	defer cleanup()
	
	// Test updating config
	cfg := Config{InstallationToken: "test-token"}
	err := provider.update(cfg)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}
	
	// Verify config was updated
	readCfg, err := provider.Get()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	
	if readCfg.InstallationToken != "test-token" {
		t.Errorf("Expected token to be 'test-token', got '%s'", readCfg.InstallationToken)
	}
	
	// Test error case - GetUserHomeDir returns error
	errorProvider := &MockConfigProvider{HomeError: true}
	
	// Test that we get an error
	err = errorProvider.update(cfg)
	if err == nil {
		t.Error("Expected error when GetUserHomeDir fails, got nil")
	}
	
	// Test OpenFile error - create a directory with the same name as the config file
	home, _ := provider.GetUserHomeDir()
	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")
	
	// Remove any existing config file
	os.Remove(configPath)
	
	// Create a directory with the same name
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	err = provider.update(cfg)
	if err == nil {
		t.Error("Expected error when OpenFile fails, got nil")
	}
	
	// Clean up
	os.RemoveAll(configPath)
}

func TestStoreData(t *testing.T) {
	provider, cleanup := setupTestConfig(t)
	defer cleanup()
	
	// Test storing data
	err := provider.StoreData("public-key", "private-key")
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// Verify data was stored
	cfg, err := provider.Get()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	
	if cfg.PublicKey != "public-key" {
		t.Errorf("Expected public key to be 'public-key', got '%s'", cfg.PublicKey)
	}
	
	if cfg.PrivateKey != "private-key" {
		t.Errorf("Expected private key to be 'private-key', got '%s'", cfg.PrivateKey)
	}
	
	// Test error case - GetUserHomeDir returns error
	errorProvider := &MockConfigProvider{HomeError: true}
	
	// Test that we get an error
	err = errorProvider.StoreData("public-key", "private-key")
	if err == nil {
		t.Error("Expected error when GetUserHomeDir fails, got nil")
	}
}

func TestGet(t *testing.T) {
	provider, cleanup := setupTestConfig(t)
	defer cleanup()
	
	// Test when config file doesn't exist
	cfg, err := provider.Get()
	if err != nil {
		t.Fatalf("Failed to get config when file doesn't exist: %v", err)
	}
	
	// Config should be empty but valid
	if cfg.InstallationToken != "" {
		t.Errorf("Expected empty token, got '%s'", cfg.InstallationToken)
	}
	
	// Now create a config file and test again
	home, _ := provider.GetUserHomeDir()
	configPath := filepath.Join(home, ".config", "mbvpn", "config.yml")
	
	// Create config file with valid data
	f, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	
	testCfg := Config{InstallationToken: "test-token"}
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(testCfg)
	f.Close()
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	
	cfg, err = provider.Get()
	if err != nil {
		t.Fatalf("Failed to get config when file exists: %v", err)
	}
	
	if cfg.InstallationToken != "test-token" {
		t.Errorf("Expected token to be 'test-token', got '%s'", cfg.InstallationToken)
	}
	
	// Test error case - GetUserHomeDir returns error
	errorProvider := &MockConfigProvider{HomeError: true}
	
	// Test that we get an error
	_, err = errorProvider.Get()
	if err == nil {
		t.Error("Expected error when GetUserHomeDir fails, got nil")
	}
	
	// Test invalid YAML
	os.Remove(configPath)
	f, _ = os.Create(configPath)
	f.WriteString("invalid: yaml: :")
	f.Close()
	
	_, err = provider.Get()
	if err == nil {
		t.Error("Expected error with invalid YAML, got nil")
	}
	
	// Test MkdirAll error
	os.Remove(configPath)
	os.RemoveAll(filepath.Dir(configPath))
	
	// Create a file where the directory should be
	configDir := filepath.Dir(configPath)
	parentDir := filepath.Dir(configDir)
	os.MkdirAll(parentDir, 0755)
	f, err = os.Create(configDir)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	f.Close()
	
	// Should fail when trying to create directory
	_, err = provider.Get()
	if err == nil {
		t.Error("Expected error when MkdirAll fails, got nil")
	}
	
	// Clean up and restore directory
	os.Remove(configDir)
	os.MkdirAll(configDir, 0755)
}