package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewConfigFileMachineIdProvider(t *testing.T) {
	provider := NewConfigFileMachineIdProvider()
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
	
	_, ok := provider.(*ConfigFileMachineIdProvider)
	if !ok {
		t.Fatal("Expected provider to be of type *ConfigFileMachineIdProvider")
	}
}

func TestConfigFileMachineIdProviderGet(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Save the original homeDir resolution function and restore it after test
	originalUserCurrent := userCurrentFunc
	defer func() { userCurrentFunc = originalUserCurrent }()
	
	// Mock the user.Current() function
	userCurrentFunc = func() (*user.User, error) {
		return &user.User{HomeDir: tempDir}, nil
	}
	
	configPath := filepath.Join(tempDir, ".config", "mbvpn")
	machineIDPath := filepath.Join(configPath, "machine-id")
	
	t.Run("creates new machine ID when file doesn't exist", func(t *testing.T) {
		// Clean up any existing file
		os.RemoveAll(configPath)
		
		provider := &ConfigFileMachineIdProvider{}
		id, err := provider.Get()
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		if id == "" {
			t.Fatal("Expected non-empty machine ID")
		}
		
		// Verify file was created
		if _, err := os.Stat(machineIDPath); os.IsNotExist(err) {
			t.Fatal("Expected machine-id file to be created")
		}
		
		// Verify ID in file matches returned ID
		data, err := os.ReadFile(machineIDPath)
		if err != nil {
			t.Fatalf("Failed to read machine-id file: %v", err)
		}
		
		if string(data) != id {
			t.Fatalf("Expected file content to match returned ID. Got: %s, want: %s", string(data), id)
		}
		
		// Verify the ID is a valid UUID
		_, err = uuid.Parse(id)
		if err != nil {
			t.Fatalf("Expected a valid UUID, got: %s", id)
		}
	})
	
	t.Run("returns existing machine ID when file exists", func(t *testing.T) {
		// Clean up and recreate the directory
		os.RemoveAll(configPath)
		os.MkdirAll(configPath, 0700)
		
		// Create a machine-id file with a known value
		expectedID := "test-machine-id"
		err := os.WriteFile(machineIDPath, []byte(expectedID), 0600)
		if err != nil {
			t.Fatalf("Failed to create test machine-id file: %v", err)
		}
		
		provider := &ConfigFileMachineIdProvider{}
		id, err := provider.Get()
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		if id != expectedID {
			t.Fatalf("Expected ID %s, got: %s", expectedID, id)
		}
	})
	
	t.Run("trims newline from machine ID", func(t *testing.T) {
		// Clean up and recreate the directory
		os.RemoveAll(configPath)
		os.MkdirAll(configPath, 0700)
		
		// Create a machine-id file with a trailing newline
		expectedID := "test-machine-id"
		err := os.WriteFile(machineIDPath, []byte(expectedID+"\n"), 0600)
		if err != nil {
			t.Fatalf("Failed to create test machine-id file: %v", err)
		}
		
		provider := &ConfigFileMachineIdProvider{}
		id, err := provider.Get()
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		if id != expectedID {
			t.Fatalf("Expected ID %s, got: %s", expectedID, id)
		}
	})
	
	t.Run("handles empty file", func(t *testing.T) {
		// Clean up and recreate the directory
		os.RemoveAll(configPath)
		os.MkdirAll(configPath, 0700)
		
		// Create an empty machine-id file
		err := os.WriteFile(machineIDPath, []byte{}, 0600)
		if err != nil {
			t.Fatalf("Failed to create empty machine-id file: %v", err)
		}
		
		provider := &ConfigFileMachineIdProvider{}
		id, err := provider.Get()
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		if id == "" {
			t.Fatal("Expected non-empty machine ID")
		}
		
		// Verify the ID is a valid UUID
		_, err = uuid.Parse(id)
		if err != nil {
			t.Fatalf("Expected a valid UUID, got: %s", id)
		}
	})
	
	t.Run("handles user.Current error", func(t *testing.T) {
		// Mock the user.Current() function to return an error
		userCurrentFunc = func() (*user.User, error) {
			return nil, fmt.Errorf("mock user.Current error")
		}
		
		provider := &ConfigFileMachineIdProvider{}
		_, err := provider.Get()
		
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		
		if err.Error() != "failed to get current user: mock user.Current error" {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})
	
	t.Run("handles MkdirAll error", func(t *testing.T) {
		// Restore the original user.Current function
		userCurrentFunc = func() (*user.User, error) {
			return &user.User{HomeDir: "/non-existent-directory-that-cannot-be-created"}, nil
		}
		
		// Create a read-only directory to cause MkdirAll to fail
		provider := &ConfigFileMachineIdProvider{}
		_, err := provider.Get()
		
		// This test is system-dependent and may not always fail on MkdirAll
		// So we'll check if the error contains our expected prefix
		if err == nil {
			t.Skip("MkdirAll did not fail, skipping test")
		}
		
		if err != nil && !strings.Contains(err.Error(), "failed to create config directory") {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})
	
	t.Run("handles WriteFile error", func(t *testing.T) {
		// Restore the user.Current function
		userCurrentFunc = func() (*user.User, error) {
			return &user.User{HomeDir: tempDir}, nil
		}
		
		// Clean up and recreate the directory, but make it read-only
		os.RemoveAll(configPath)
		os.MkdirAll(configPath, 0500) // read + execute, but no write permission
		
		// This test may only work when running as a non-root user
		// as root can write to read-only directories
		provider := &ConfigFileMachineIdProvider{}
		_, err := provider.Get()
		
		// If we're running as root, the test might not fail
		if err == nil {
			if os.Geteuid() == 0 {
				t.Skip("Test running as root, skipping WriteFile error test")
			} else {
				t.Fatal("Expected WriteFile to fail with read-only directory")
			}
		}
		
		if err != nil && !strings.Contains(err.Error(), "failed to write machine-id file") {
			t.Fatalf("Unexpected error message: %v", err)
		}
		
		// Reset permissions for cleanup
		os.Chmod(configPath, 0700)
	})
}