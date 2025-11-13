package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultDirectoryProvider_GetConfigDir(t *testing.T) {
	provider := NewDefaultDirectoryProvider()

	configDir, err := provider.GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir should succeed: %v", err)
	}

	if configDir == "" {
		t.Error("GetConfigDir should return non-empty directory")
	}

	if !strings.HasSuffix(configDir, "mbvpn") {
		t.Errorf("GetConfigDir should end with 'mbvpn', got: %s", configDir)
	}
}

func TestDefaultDirectoryProvider_GetServersDir(t *testing.T) {
	provider := NewDefaultDirectoryProvider()

	serversDir, err := provider.GetServersDir()
	if err != nil {
		t.Fatalf("GetServersDir should succeed: %v", err)
	}

	if !strings.HasSuffix(serversDir, filepath.Join("mbvpn", "servers")) {
		t.Errorf("GetServersDir should end with 'mbvpn/servers', got: %s", serversDir)
	}
}

func TestDefaultDirectoryProvider_GetConfigFile(t *testing.T) {
	provider := NewDefaultDirectoryProvider()

	configFile, err := provider.GetConfigFile()
	if err != nil {
		t.Fatalf("GetConfigFile should succeed: %v", err)
	}

	if !strings.HasSuffix(configFile, filepath.Join("mbvpn", "config.yml")) {
		t.Errorf("GetConfigFile should end with 'mbvpn/config.yml', got: %s", configFile)
	}
}

func TestDefaultDirectoryProvider_GetMachineIDFile(t *testing.T) {
	provider := NewDefaultDirectoryProvider()

	machineIDFile, err := provider.GetMachineIDFile()
	if err != nil {
		t.Fatalf("GetMachineIDFile should succeed: %v", err)
	}

	if !strings.HasSuffix(machineIDFile, filepath.Join("mbvpn", "machine-id")) {
		t.Errorf("GetMachineIDFile should end with 'mbvpn/machine-id', got: %s", machineIDFile)
	}
}

func TestDefaultDirectoryProvider_GetServersFile(t *testing.T) {
	provider := NewDefaultDirectoryProvider()

	serversFile, err := provider.GetServersFile()
	if err != nil {
		t.Fatalf("GetServersFile should succeed: %v", err)
	}

	if !strings.HasSuffix(serversFile, filepath.Join("mbvpn", "servers.json")) {
		t.Errorf("GetServersFile should end with 'mbvpn/servers.json', got: %s", serversFile)
	}
}

func TestDefaultDirectoryProvider_ErrorCase(t *testing.T) {
	// Test error case by unsetting both XDG_CONFIG_HOME and HOME
	oldConfigHome := os.Getenv("XDG_CONFIG_HOME")
	oldHome := os.Getenv("HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Unsetenv("HOME")
	defer func() {
		if oldConfigHome != "" {
			os.Setenv("XDG_CONFIG_HOME", oldConfigHome)
		}
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		}
	}()

	provider := NewDefaultDirectoryProvider()

	_, err := provider.GetConfigDir()
	if err == nil {
		t.Error("GetConfigDir should fail when both XDG_CONFIG_HOME and HOME are unset")
	}

	if !strings.Contains(err.Error(), "failed to get user config directory") {
		t.Errorf("Expected error about config directory, got: %v", err)
	}
}

func TestTestDirectoryProvider_GetConfigDir(t *testing.T) {
	tempDir := t.TempDir()
	provider := NewTestDirectoryProvider(tempDir)

	configDir, err := provider.GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir should succeed: %v", err)
	}

	expectedDir := filepath.Join(tempDir, "mbvpn")
	if configDir != expectedDir {
		t.Errorf("Expected config dir '%s', got '%s'", expectedDir, configDir)
	}
}

func TestTestDirectoryProvider_GetServersDir(t *testing.T) {
	tempDir := t.TempDir()
	provider := NewTestDirectoryProvider(tempDir)

	serversDir, err := provider.GetServersDir()
	if err != nil {
		t.Fatalf("GetServersDir should succeed: %v", err)
	}

	expectedDir := filepath.Join(tempDir, "mbvpn", "servers")
	if serversDir != expectedDir {
		t.Errorf("Expected servers dir '%s', got '%s'", expectedDir, serversDir)
	}
}

func TestTestDirectoryProvider_GetConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	provider := NewTestDirectoryProvider(tempDir)

	configFile, err := provider.GetConfigFile()
	if err != nil {
		t.Fatalf("GetConfigFile should succeed: %v", err)
	}

	expectedFile := filepath.Join(tempDir, "mbvpn", "config.yml")
	if configFile != expectedFile {
		t.Errorf("Expected config file '%s', got '%s'", expectedFile, configFile)
	}
}

func TestTestDirectoryProvider_GetMachineIDFile(t *testing.T) {
	tempDir := t.TempDir()
	provider := NewTestDirectoryProvider(tempDir)

	machineIDFile, err := provider.GetMachineIDFile()
	if err != nil {
		t.Fatalf("GetMachineIDFile should succeed: %v", err)
	}

	expectedFile := filepath.Join(tempDir, "mbvpn", "machine-id")
	if machineIDFile != expectedFile {
		t.Errorf("Expected machine ID file '%s', got '%s'", expectedFile, machineIDFile)
	}
}

func TestTestDirectoryProvider_GetServersFile(t *testing.T) {
	tempDir := t.TempDir()
	provider := NewTestDirectoryProvider(tempDir)

	serversFile, err := provider.GetServersFile()
	if err != nil {
		t.Fatalf("GetServersFile should succeed: %v", err)
	}

	expectedFile := filepath.Join(tempDir, "mbvpn", "servers.json")
	if serversFile != expectedFile {
		t.Errorf("Expected servers file '%s', got '%s'", expectedFile, serversFile)
	}
}

func TestTestDirectoryProvider_Integration(t *testing.T) {
	tempDir := t.TempDir()
	provider := NewTestDirectoryProvider(tempDir)

	// Test all methods return paths under the same base directory
	configDir, _ := provider.GetConfigDir()
	serversDir, _ := provider.GetServersDir()
	configFile, _ := provider.GetConfigFile()
	machineIDFile, _ := provider.GetMachineIDFile()
	serversFile, _ := provider.GetServersFile()

	paths := []string{configDir, serversDir, configFile, machineIDFile, serversFile}
	for _, path := range paths {
		if !strings.HasPrefix(path, tempDir) {
			t.Errorf("Path should be under temp directory %s, got: %s", tempDir, path)
		}
	}
}

func TestDirectoryProvider_Interface(t *testing.T) {
	// Verify both implementations satisfy the interface
	var _ DirectoryProvider = &DefaultDirectoryProvider{}
	var _ DirectoryProvider = &TestDirectoryProvider{}
}
