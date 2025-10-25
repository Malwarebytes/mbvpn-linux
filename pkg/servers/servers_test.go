package servers

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Malwarebytes/mbvpn-linux/pkg/remote"
)

// setupTestDir creates a temporary directory for testing
func setupTestDir(t *testing.T) (string, func()) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "mbvpn")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test config directory: %v", err)
	}

	// Set HOME to temp directory for tests
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)

	return tempDir, func() {
		os.Setenv("HOME", oldHome)
	}
}

// createTestLocations creates sample VPN locations for testing
func createTestLocations() *remote.VpnLocations {
	return &remote.VpnLocations{
		Countries: []remote.Country{
			{
				Code:        "US",
				Name:        "United States",
				DefaultCity: "NYC",
				Cities: []remote.City{
					{
						Code:      "NYC",
						UnLocode:  "USNYC",
						Name:      "New York",
						Latitude:  40.7128,
						Longitude: -74.0060,
						Servers: []remote.Server{
							{
								Hostname:     "vpn01.us.example.com",
								IPv4AddrIn:   "203.0.113.1",
								IPv6AddrIn:   "2001:db8::1",
								Weight:       100,
								PublicKey:    "test-public-key-1",
								PortRanges:   []remote.PortRange{{From: 51820, To: 51830}},
								MultihopPort: 51821,
								IPv4Gateway:  "203.0.113.254",
								IPv6Gateway:  "2001:db8::254",
							},
							{
								Hostname:  "vpn02.us.example.com",
								Weight:    80,
								PublicKey: "test-public-key-2",
							},
						},
					},
					{
						Code: "LAX",
						Name: "Los Angeles",
						Servers: []remote.Server{
							{
								Hostname:  "vpn01.lax.example.com",
								Weight:    90,
								PublicKey: "test-public-key-3",
							},
						},
					},
				},
			},
			{
				Code: "UK",
				Name: "United Kingdom",
				Cities: []remote.City{
					{
						Code: "LON",
						Name: "London",
						Servers: []remote.Server{
							{
								Hostname:  "vpn01.uk.example.com",
								Weight:    95,
								PublicKey: "test-public-key-4",
							},
						},
					},
				},
			},
		},
	}
}

func TestNewDefaultServerStorage(t *testing.T) {
	storage := NewDefaultServerStorage()
	if storage == nil {
		t.Fatal("NewDefaultServerStorage should return a non-nil storage")
	}

	// Verify it implements the interface
	var _ ServerStorage = storage
}

func TestDefaultServerStorage_Save(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	storage := &DefaultServerStorage{}
	locations := createTestLocations()

	t.Run("Save locations successfully", func(t *testing.T) {
		err := storage.Save(locations)
		if err != nil {
			t.Fatalf("Save should succeed: %v", err)
		}

		// Verify file was created
		home, _ := os.UserHomeDir()
		filePath := filepath.Join(home, ".config", "mbvpn", "servers.json")

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Error("servers.json file should be created")
		}

		// Verify file contents
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		var saved remote.VpnLocations
		err = json.Unmarshal(data, &saved)
		if err != nil {
			t.Fatalf("Failed to unmarshal saved data: %v", err)
		}

		if len(saved.Countries) != len(locations.Countries) {
			t.Errorf("Expected %d countries, got %d", len(locations.Countries), len(saved.Countries))
		}
	})

	t.Run("Save with nil locations", func(t *testing.T) {
		err := storage.Save(nil)
		if err != nil {
			t.Fatalf("Save should handle nil locations: %v", err)
		}
	})
}

func TestDefaultServerStorage_Save_Errors(t *testing.T) {
	storage := &DefaultServerStorage{}
	locations := createTestLocations()

	t.Run("Error when HOME not set", func(t *testing.T) {
		oldHome := os.Getenv("HOME")
		os.Unsetenv("HOME")
		defer os.Setenv("HOME", oldHome)

		err := storage.Save(locations)
		if err == nil {
			t.Error("Save should fail when HOME is not set")
		}
		if !strings.Contains(err.Error(), "failed to get user home directory") {
			t.Errorf("Expected error about home directory, got: %v", err)
		}
	})

	t.Run("Error when directory is not writable", func(t *testing.T) {
		tempDir, cleanup := setupTestDir(t)
		defer cleanup()

		// Make config directory read-only
		configDir := filepath.Join(tempDir, ".config", "mbvpn")
		err := os.Chmod(configDir, 0555) // read and execute only
		if err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		// Ensure permissions cleanup happens even if test fails
		defer func() {
			os.Chmod(configDir, 0755)
		}()

		err = storage.Save(locations)
		// Running as root may bypass permission restrictions
		if os.Getuid() == 0 {
			t.Skip("Test running as root, skipping permission test")
			return
		}
		if err == nil {
			t.Error("Save should fail when directory is not writable")
			return
		}
		if !strings.Contains(err.Error(), "failed to create file") {
			t.Errorf("Expected error about creating file, got: %v", err)
		}
	})
}

func TestDefaultServerStorage_Get(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	storage := &DefaultServerStorage{}
	originalLocations := createTestLocations()

	t.Run("Get after save", func(t *testing.T) {
		// First save the data
		err := storage.Save(originalLocations)
		if err != nil {
			t.Fatalf("Failed to save locations: %v", err)
		}

		// Then retrieve it
		locations, err := storage.Get()
		if err != nil {
			t.Fatalf("Get should succeed: %v", err)
		}

		if locations == nil {
			t.Fatal("Get should return non-nil locations")
		}

		if len(locations.Countries) != len(originalLocations.Countries) {
			t.Errorf("Expected %d countries, got %d", len(originalLocations.Countries), len(locations.Countries))
		}

		// Verify specific data
		if locations.Countries[0].Code != "US" {
			t.Errorf("Expected first country code 'US', got '%s'", locations.Countries[0].Code)
		}
	})

	t.Run("Get non-existent file", func(t *testing.T) {
		// Create a fresh temporary directory for this test
		freshTempDir, freshCleanup := setupTestDir(t)
		defer freshCleanup()

		// Use a fresh storage instance with no saved data
		newStorage := &DefaultServerStorage{}

		_, err := newStorage.Get()
		if err == nil {
			t.Error("Get should fail when file doesn't exist")
		} else if !strings.Contains(err.Error(), "failed to open file") {
			t.Errorf("Expected error about opening file, got: %v", err)
		}

		// Suppress unused variable warning
		_ = freshTempDir
	})

	t.Run("Get corrupted file", func(t *testing.T) {
		// Create corrupted JSON file
		home := tempDir
		filePath := filepath.Join(home, ".config", "mbvpn", "servers.json")

		err := os.WriteFile(filePath, []byte("invalid json {"), 0644)
		if err != nil {
			t.Fatalf("Failed to create corrupted file: %v", err)
		}

		_, err = storage.Get()
		if err == nil {
			t.Error("Get should fail with corrupted JSON")
		}
		if !strings.Contains(err.Error(), "failed to decode") {
			t.Errorf("Expected decode error, got: %v", err)
		}
	})
}

func TestDefaultServerStorage_Get_Errors(t *testing.T) {
	storage := &DefaultServerStorage{}

	t.Run("Error when HOME not set", func(t *testing.T) {
		oldHome := os.Getenv("HOME")
		os.Unsetenv("HOME")
		defer os.Setenv("HOME", oldHome)

		_, err := storage.Get()
		if err == nil {
			t.Error("Get should fail when HOME is not set")
		}
		if !strings.Contains(err.Error(), "failed to get user home directory") {
			t.Errorf("Expected error about home directory, got: %v", err)
		}
	})
}

func TestDefaultServerStorage_GetByServerName(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	storage := &DefaultServerStorage{}
	locations := createTestLocations()

	// Save test data
	err := storage.Save(locations)
	if err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	t.Run("Find by full hostname", func(t *testing.T) {
		server, err := storage.GetByServerName("vpn01.us.example.com")
		if err != nil {
			t.Fatalf("GetByServerName should succeed: %v", err)
		}
		if server == nil {
			t.Fatal("Should find server by full hostname")
		}
		if server.Hostname != "vpn01.us.example.com" {
			t.Errorf("Expected hostname 'vpn01.us.example.com', got '%s'", server.Hostname)
		}
	})

	t.Run("Find by short hostname", func(t *testing.T) {
		server, err := storage.GetByServerName("vpn01")
		if err != nil {
			t.Fatalf("GetByServerName should succeed: %v", err)
		}
		if server == nil {
			t.Fatal("Should find server by short hostname")
		}
		// Should find the first server with hostname starting with "vpn01"
		if !strings.HasPrefix(server.Hostname, "vpn01") {
			t.Errorf("Expected hostname starting with 'vpn01', got '%s'", server.Hostname)
		}
	})

	t.Run("Find by city code", func(t *testing.T) {
		server, err := storage.GetByServerName("NYC")
		if err != nil {
			t.Fatalf("GetByServerName should succeed: %v", err)
		}
		if server == nil {
			t.Fatal("Should find server by city code")
		}
		// Should return a random server from NYC
	})

	t.Run("Find by city name", func(t *testing.T) {
		server, err := storage.GetByServerName("New York")
		if err != nil {
			t.Fatalf("GetByServerName should succeed: %v", err)
		}
		if server == nil {
			t.Fatal("Should find server by city name")
		}
	})

	t.Run("Find by country code", func(t *testing.T) {
		server, err := storage.GetByServerName("US")
		if err != nil {
			t.Fatalf("GetByServerName should succeed: %v", err)
		}
		if server == nil {
			t.Fatal("Should find server by country code")
		}
	})

	t.Run("Find by country name", func(t *testing.T) {
		server, err := storage.GetByServerName("United States")
		if err != nil {
			t.Fatalf("GetByServerName should succeed: %v", err)
		}
		if server == nil {
			t.Fatal("Should find server by country name")
		}
	})

	t.Run("Server not found", func(t *testing.T) {
		server, err := storage.GetByServerName("nonexistent")
		if err != nil {
			t.Fatalf("GetByServerName should not error for non-existent server: %v", err)
		}
		if server != nil {
			t.Error("Should return nil for non-existent server")
		}
	})

	t.Run("Empty server name", func(t *testing.T) {
		server, err := storage.GetByServerName("")
		if err != nil {
			t.Fatalf("GetByServerName should not error for empty name: %v", err)
		}
		if server != nil {
			t.Error("Should return nil for empty server name")
		}
	})
}

func TestDefaultServerStorage_GetByServerName_Errors(t *testing.T) {
	storage := &DefaultServerStorage{}

	t.Run("Error when no data saved", func(t *testing.T) {
		_, cleanup := setupTestDir(t)
		defer cleanup()

		_, err := storage.GetByServerName("test")
		if err == nil {
			t.Error("GetByServerName should fail when no data is saved")
		}
		if !strings.Contains(err.Error(), "failed to get locations") {
			t.Errorf("Expected error about getting locations, got: %v", err)
		}
	})
}

func TestRandomInt(t *testing.T) {
	// Set a seed for reproducible tests
	rand.Seed(time.Now().UnixNano())

	t.Run("Random int within range", func(t *testing.T) {
		max := 10
		for i := 0; i < 100; i++ {
			result := randomInt(max)
			if result < 0 || result >= max {
				t.Errorf("randomInt(%d) returned %d, should be 0 <= result < %d", max, result, max)
			}
		}
	})

	t.Run("Random int with max 1", func(t *testing.T) {
		result := randomInt(1)
		if result != 0 {
			t.Errorf("randomInt(1) should always return 0, got %d", result)
		}
	})

	t.Run("Max zero returns zero", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("randomInt(0) panicked: %v", r)
			}
		}()

		result := randomInt(0)
		if result != 0 {
			t.Errorf("randomInt(0) should return 0, got %d", result)
		}
	})

	t.Run("Random distribution", func(t *testing.T) {
		// Test that randomInt produces different values over multiple calls
		max := 10
		results := make(map[int]bool)

		for i := 0; i < 100; i++ {
			result := randomInt(max)
			results[result] = true
		}

		// We should see at least a few different values in 100 calls
		if len(results) < 3 {
			t.Errorf("Expected at least 3 different values in 100 calls, got %d", len(results))
		}
	})
}

func TestServerStorage_Interface(t *testing.T) {
	// Test that DefaultServerStorage implements ServerStorage interface
	var storage ServerStorage = &DefaultServerStorage{}

	// Interface compliance test - this will fail to compile if interface is not implemented
	_ = storage
}

func TestDefaultServerStorage_Integration(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	storage := &DefaultServerStorage{}
	locations := createTestLocations()

	// Full integration test: save, get, and search
	t.Run("Full workflow", func(t *testing.T) {
		// 1. Save locations
		err := storage.Save(locations)
		if err != nil {
			t.Fatalf("Failed to save: %v", err)
		}

		// 2. Verify we can get them back
		retrieved, err := storage.Get()
		if err != nil {
			t.Fatalf("Failed to get: %v", err)
		}

		// 3. Verify we can search for servers
		server, err := storage.GetByServerName("vpn01.us.example.com")
		if err != nil {
			t.Fatalf("Failed to search: %v", err)
		}
		if server == nil {
			t.Fatal("Should find the server")
		}

		// 4. Verify data integrity
		if len(retrieved.Countries) != 2 {
			t.Errorf("Expected 2 countries, got %d", len(retrieved.Countries))
		}
		if server.PublicKey != "test-public-key-1" {
			t.Errorf("Expected public key 'test-public-key-1', got '%s'", server.PublicKey)
		}
	})
}

// Benchmark tests for performance
func BenchmarkDefaultServerStorage_Save(b *testing.B) {
	tempDir := b.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	configDir := filepath.Join(tempDir, ".config", "mbvpn")
	os.MkdirAll(configDir, 0755)

	storage := &DefaultServerStorage{}
	locations := createTestLocations()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Save(locations)
	}
}

func BenchmarkDefaultServerStorage_Get(b *testing.B) {
	tempDir := b.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	configDir := filepath.Join(tempDir, ".config", "mbvpn")
	os.MkdirAll(configDir, 0755)

	storage := &DefaultServerStorage{}
	locations := createTestLocations()
	storage.Save(locations) // Setup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Get()
	}
}

func BenchmarkDefaultServerStorage_GetByServerName(b *testing.B) {
	tempDir := b.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	configDir := filepath.Join(tempDir, ".config", "mbvpn")
	os.MkdirAll(configDir, 0755)

	storage := &DefaultServerStorage{}
	locations := createTestLocations()
	storage.Save(locations) // Setup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.GetByServerName("vpn01.us.example.com")
	}
}
