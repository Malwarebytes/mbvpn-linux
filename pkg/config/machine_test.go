package config

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewConfigFileMachineIdProvider(t *testing.T) {
	dirProvider := NewDefaultDirectoryProvider()
	provider := NewConfigFileMachineIdProvider(dirProvider)
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}

	_, ok := provider.(*ConfigFileMachineIdProvider)
	if !ok {
		t.Fatal("Expected provider to be of type *ConfigFileMachineIdProvider")
	}
}

func TestConfigFileMachineIdProviderGet(t *testing.T) {
	// Test that the provider can successfully generate and retrieve a machine ID
	// Now uses a temporary directory instead of the actual config directory
	tempDir := t.TempDir()
	dirProvider := NewTestDirectoryProvider(tempDir)
	provider := NewConfigFileMachineIdProvider(dirProvider)

	// First call - should create new ID
	id1, err := provider.Get()
	if err != nil {
		t.Fatalf("Expected no error on first call, got: %v", err)
	}

	if id1 == "" {
		t.Fatal("Expected non-empty machine ID")
	}

	// Verify the ID is a valid UUID
	_, uuidErr := uuid.Parse(id1)
	if uuidErr != nil {
		t.Fatalf("Machine ID should be a valid UUID, got: %s, error: %v", id1, uuidErr)
	}

	// Second call - should return the same ID
	id2, err := provider.Get()
	if err != nil {
		t.Fatalf("Expected no error on second call, got: %v", err)
	}

	if id1 != id2 {
		t.Fatalf("Expected consistent machine ID. First: %s, Second: %s", id1, id2)
	}
}

func TestConfigFileMachineIdProviderGet_ErrorCase(t *testing.T) {
	// Test error case using MockDirectoryProvider
	mockDirProvider := &MockDirectoryProvider{ConfigDirError: true}
	provider := NewConfigFileMachineIdProvider(mockDirProvider)

	_, err := provider.Get()

	if err == nil {
		t.Fatal("Expected an error when directory provider fails, got nil")
	}
}
