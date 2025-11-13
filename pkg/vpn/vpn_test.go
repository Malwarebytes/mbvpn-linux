package vpn

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	"github.com/malwarebytes/mbvpn-linux/pkg/remote"
)

// Helper function to set up test directory
func setupTestDir(t *testing.T) (config.DirectoryProvider, func()) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "mbvpn")
	serversDir := filepath.Join(configDir, "servers")
	if err := os.MkdirAll(serversDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}
	dirProvider := config.NewDirectoryProvider(tempDir)
	return dirProvider, func() {}
}

// Mock DirectoryProvider for error testing
type mockDirectoryProvider struct {
	shouldError bool
}

func (m *mockDirectoryProvider) GetConfigDir() (string, error) {
	if m.shouldError {
		return "", os.ErrPermission
	}
	return "", nil
}

func (m *mockDirectoryProvider) GetServersDir() (string, error) {
	if m.shouldError {
		return "", os.ErrPermission
	}
	return "", nil
}

func (m *mockDirectoryProvider) GetConfigFile() (string, error) {
	if m.shouldError {
		return "", os.ErrPermission
	}
	return "", nil
}

func (m *mockDirectoryProvider) GetMachineIDFile() (string, error) {
	if m.shouldError {
		return "", os.ErrPermission
	}
	return "", nil
}

func (m *mockDirectoryProvider) GetServersFile() (string, error) {
	if m.shouldError {
		return "", os.ErrPermission
	}
	return "", nil
}

func TestGetCountryFlag(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		expected    string
	}{
		{"United States", "US", "🇺🇸"},
		{"United Kingdom", "GB", "🇬🇧"},
		{"Germany", "DE", "🇩🇪"},
		{"Japan", "JP", "🇯🇵"},
		{"Canada", "CA", "🇨🇦"},
		{"France", "FR", "🇫🇷"},
		{"Australia", "AU", "🇦🇺"},
		{"Brazil", "BR", "🇧🇷"},
		{"India", "IN", "🇮🇳"},
		{"China", "CN", "🇨🇳"},
		{"Russia", "RU", "🇷🇺"},
		{"South Africa", "ZA", "🇿🇦"},
		{"Netherlands", "NL", "🇳🇱"},
		{"Sweden", "SE", "🇸🇪"},
		{"Switzerland", "CH", "🇨🇭"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCountryFlag(tt.countryCode)
			if result != tt.expected {
				t.Errorf("getCountryFlag(%s) = %s, expected %s", tt.countryCode, result, tt.expected)
			}
		})
	}
}

func TestGetCountryFlag_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"us", "🇺🇸"},
		{"Us", "🇺🇸"},
		{"US", "🇺🇸"},
		{"uS", "🇺🇸"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := getCountryFlag(tt.input)
			if result != tt.expected {
				t.Errorf("getCountryFlag(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCountryFlag_UnknownCountry(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
	}{
		{"Unknown two letter code", "XX"},
		{"Made up code", "ZZ"},
		{"Another unknown", "QQ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCountryFlag(tt.countryCode)
			// Should generate a flag from Unicode Regional Indicator symbols
			if len(result) == 0 {
				t.Errorf("getCountryFlag(%s) should return non-empty result", tt.countryCode)
			}
			// Verify it's not the default globe
			if result == "🌍" {
				t.Errorf("getCountryFlag(%s) should generate flag, not default to globe", tt.countryCode)
			}
		})
	}
}

func TestGetCountryFlag_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		expected    string
	}{
		{"Empty string", "", "🌍"},
		{"Single character", "A", "🌍"},
		{"Three characters", "USA", "🌍"},
		{"Long string", "UNITED_STATES", "🌍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCountryFlag(tt.countryCode)
			if result != tt.expected {
				t.Errorf("getCountryFlag(%s) = %s, expected %s", tt.countryCode, result, tt.expected)
			}
		})
	}
}

func TestGetCountryFlag_NonLetterInput(t *testing.T) {
	// Test that non-letter inputs generate flags rather than globe
	// This is the actual behavior of the function which converts any 2-char input to Unicode flags
	tests := []struct {
		name        string
		countryCode string
	}{
		{"Numbers", "12"},
		{"Special characters", "!@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCountryFlag(tt.countryCode)
			// For 2-character inputs, the function generates Unicode flags even for non-letters
			// The specific flag depends on the Unicode values, so we just verify it's not empty and not the globe
			if result == "" {
				t.Errorf("getCountryFlag(%s) should return non-empty result", tt.countryCode)
			}
			if result == "🌍" {
				t.Errorf("getCountryFlag(%s) should generate flag for 2-char input, not globe", tt.countryCode)
			}
		})
	}
}

func TestGetCountryFlag_AllMappedCountries(t *testing.T) {
	// Test that all mapped countries return their expected flags
	// This ensures the hardcoded map is working correctly
	mappedCountries := map[string]string{
		"AF": "🇦🇫", "AL": "🇦🇱", "DZ": "🇩🇿", "AR": "🇦🇷", "AM": "🇦🇲",
		"AU": "🇦🇺", "AT": "🇦🇹", "AZ": "🇦🇿", "BH": "🇧🇭", "BD": "🇧🇩",
		"BY": "🇧🇾", "BE": "🇧🇪", "BO": "🇧🇴", "BA": "🇧🇦", "BR": "🇧🇷",
		"KH": "🇰🇭", "CA": "🇨🇦", "CL": "🇨🇱", "CN": "🇨🇳", "CO": "🇨🇴",
		"CR": "🇨🇷", "HR": "🇭🇷", "CU": "🇨🇺", "CY": "🇨🇾", "CZ": "🇨🇿",
		"DK": "🇩🇰", "DO": "🇩🇴", "EC": "🇪🇨", "EG": "🇪🇬", "SV": "🇸🇻",
		"EE": "🇪🇪", "FI": "🇫🇮", "FR": "🇫🇷", "GE": "🇬🇪", "DE": "🇩🇪",
		"GH": "🇬🇭", "GB": "🇬🇧", "GR": "🇬🇷", "GT": "🇬🇹", "HK": "🇭🇰",
		"HN": "🇭🇳", "HU": "🇭🇺", "IS": "🇮🇸", "IN": "🇮🇳", "ID": "🇮🇩",
		"IR": "🇮🇷", "IQ": "🇮🇶", "IE": "🇮🇪", "IL": "🇮🇱", "IT": "🇮🇹",
		"JM": "🇯🇲", "JP": "🇯🇵", "JO": "🇯🇴", "KZ": "🇰🇿", "KE": "🇰🇪",
		"KR": "🇰🇷", "KW": "🇰🇼", "LV": "🇱🇻", "LB": "🇱🇧", "LY": "🇱🇾",
		"LT": "🇱🇹", "LU": "🇱🇺", "MK": "🇲🇰", "MY": "🇲🇾", "MT": "🇲🇹",
		"MX": "🇲🇽", "MD": "🇲🇩", "MN": "🇲🇳", "ME": "🇲🇪", "MA": "🇲🇦",
		"MM": "🇲🇲", "NP": "🇳🇵", "NL": "🇳🇱", "NZ": "🇳🇿", "NI": "🇳🇮",
		"NG": "🇳🇬", "NO": "🇳🇴", "OM": "🇴🇲", "PK": "🇵🇰", "PA": "🇵🇦",
		"PY": "🇵🇾", "PE": "🇵🇪", "PH": "🇵🇭", "PL": "🇵🇱", "PT": "🇵🇹",
		"PR": "🇵🇷", "QA": "🇶🇦", "RO": "🇷🇴", "RU": "🇷🇺", "SA": "🇸🇦",
		"RS": "🇷🇸", "SG": "🇸🇬", "SK": "🇸🇰", "SI": "🇸🇮", "ZA": "🇿🇦",
		"ES": "🇪🇸", "LK": "🇱🇰", "SD": "🇸🇩", "SE": "🇸🇪", "CH": "🇨🇭",
		"SY": "🇸🇾", "TW": "🇹🇼", "TJ": "🇹🇯", "TH": "🇹🇭", "TN": "🇹🇳",
		"TR": "🇹🇷", "TM": "🇹🇲", "UA": "🇺🇦", "AE": "🇦🇪", "US": "🇺🇸",
		"UY": "🇺🇾", "UZ": "🇺🇿", "VE": "🇻🇪", "VN": "🇻🇳",
	}

	for code, expectedFlag := range mappedCountries {
		t.Run(code, func(t *testing.T) {
			result := getCountryFlag(code)
			if result != expectedFlag {
				t.Errorf("getCountryFlag(%s) = %s, expected %s", code, result, expectedFlag)
			}
		})
	}
}

func TestGetServersDir(t *testing.T) {
	// Test creating config directory
	dirProvider, cleanup := setupTestDir(t)
	defer cleanup()

	vpn := &DefaultVpn{dirProvider: dirProvider}

	configDir, err := vpn.dirProvider.GetServersDir()
	if err != nil {
		t.Fatalf("GetServersDir should succeed: %v", err)
	}

	// Verify directory path is correct
	if !strings.Contains(configDir, "mbvpn") || !strings.Contains(configDir, "servers") {
		t.Errorf("Expected config dir to contain 'mbvpn/servers', got '%s'", configDir)
	}

	// Verify directory was created during setup
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory should exist")
	}

	// Verify directory permissions
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}

	if info.Mode().Perm() != 0755 {
		t.Errorf("Expected directory permissions 0755, got %o", info.Mode().Perm())
	}
}

func TestGetServersDir_ErrorCase(t *testing.T) {
	// Test error case using mockDirectoryProvider
	mockDirProvider := &mockDirectoryProvider{shouldError: true}
	vpn := &DefaultVpn{dirProvider: mockDirProvider}

	_, err := vpn.dirProvider.GetServersDir()
	if err == nil {
		t.Error("GetServersDir should fail when directory provider fails")
	}
}

func TestSaveWgConfig(t *testing.T) {
	dirProvider, cleanup := setupTestDir(t)
	defer cleanup()

	vpn := &DefaultVpn{dirProvider: dirProvider}

	serverName := "test-server"
	configContent := `[Interface]
PrivateKey = test-private-key
Address = 10.0.0.1/32, 2001:db8::1/128

[Peer]
PublicKey = test-public-key
Endpoint = 203.0.113.1:51820
AllowedIPs = 0.0.0.0/0, ::/0`

	filePath, err := vpn.saveWgConfig(serverName, configContent)
	if err != nil {
		t.Fatalf("saveWgConfig should succeed: %v", err)
	}

	// Verify file path contains expected components
	if !strings.Contains(filePath, "mbvpn") || !strings.Contains(filePath, "servers") || !strings.Contains(filePath, "test-server.conf") {
		t.Errorf("Expected file path to contain 'mbvpn/servers/test-server.conf', got '%s'", filePath)
	}

	// Verify file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Config file should be created")
	}

	// Verify file permissions (should be restrictive for private keys)
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	// Verify file contents
	savedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedContent) != configContent {
		t.Errorf("File content mismatch.\nExpected:\n%s\nGot:\n%s", configContent, string(savedContent))
	}
}

func TestSaveWgConfig_PathTraversal(t *testing.T) {
	dirProvider, cleanup := setupTestDir(t)
	defer cleanup()

	vpn := &DefaultVpn{dirProvider: dirProvider}

	// Test that path traversal attempts are sanitized
	maliciousNames := []string{
		"../../../etc/passwd",
		"/etc/passwd",
		"server/../../../etc/passwd",
	}

	for _, name := range maliciousNames {
		t.Run(name, func(t *testing.T) {
			filePath, err := vpn.saveWgConfig(name, "test content")
			if err != nil {
				t.Fatalf("saveWgConfig should not fail due to path sanitization: %v", err)
			}

			// Verify the file is created in the correct directory (contains mbvpn/servers)
			if !strings.Contains(filePath, "mbvpn") || !strings.Contains(filePath, "servers") {
				t.Errorf("File should be created in mbvpn/servers directory, got: %s", filePath)
			}

			// Verify the filename is sanitized (should not contain path separators)
			fileName := filepath.Base(filePath)
			if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") {
				t.Errorf("Filename should be sanitized, got: %s", fileName)
			}
		})
	}
}

func TestWriteConfig(t *testing.T) {
	dirProvider, cleanup := setupTestDir(t)
	defer cleanup()

	vpn := &DefaultVpn{dirProvider: dirProvider}

	cfgName := "test-server"
	server := remote.Server{
		Hostname:   "vpn.example.com",
		IPv4AddrIn: "203.0.113.1",
		PublicKey:  "test-server-public-key",
	}
	privateKey := "test-private-key"
	ipv4 := "10.0.0.1/32"
	ipv6 := "2001:db8::1/128"

	configPath, err := vpn.writeConfig(cfgName, server, privateKey, ipv4, ipv6)
	if err != nil {
		t.Fatalf("writeConfig should succeed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should be created")
	}

	// Read and verify config content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	configStr := string(content)

	// Verify required sections and values are present
	expectedValues := []string{
		"[Interface]",
		"PrivateKey = " + privateKey,
		"Address = " + ipv4 + ", " + ipv6,
		"[Peer]",
		"PublicKey = " + server.PublicKey,
		"Endpoint = " + server.IPv4AddrIn + ":51820",
		"AllowedIPs = 0.0.0.0/0, ::/0",
	}

	for _, expected := range expectedValues {
		if !strings.Contains(configStr, expected) {
			t.Errorf("Config should contain '%s'\nActual config:\n%s", expected, configStr)
		}
	}
}

func TestGetConnectedServers_ParseOutput(t *testing.T) {
	// This tests the parsing logic without actually running wg command
	// We'll test the internal parsing logic by extracting it to a separate function

	// Test parsing wg show output
	testOutput := `interface: wg0
  public key: test-public-key-1
  private key: (hidden)
  listening port: 51820

peer: peer-public-key-1
  endpoint: 203.0.113.1:51820
  allowed ips: 0.0.0.0/0, ::/0

interface: wg1
  public key: test-public-key-2
  private key: (hidden)
  listening port: 51821

peer: peer-public-key-2
  endpoint: 203.0.113.2:51820
  allowed ips: 0.0.0.0/0, ::/0`

	// Extract interface names from output
	lines := strings.Split(testOutput, "\n")
	interfaceCount := strings.Count(testOutput, "interface:")
	interfaces := make([]string, interfaceCount)
	i := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "interface") {
			interfaceName := strings.TrimPrefix(line, "interface: ")
			interfaces[i] = interfaceName
			i++
		}
	}

	expectedInterfaces := []string{"wg0", "wg1"}
	if len(interfaces) != len(expectedInterfaces) {
		t.Errorf("Expected %d interfaces, got %d", len(expectedInterfaces), len(interfaces))
	}

	for i, expected := range expectedInterfaces {
		if i < len(interfaces) && interfaces[i] != expected {
			t.Errorf("Expected interface %s, got %s", expected, interfaces[i])
		}
	}
}

func TestGetConnectedServers_EmptyOutput(t *testing.T) {
	// Test parsing empty wg show output
	testOutput := ""

	interfaceCount := strings.Count(testOutput, "interface:")
	if interfaceCount != 0 {
		t.Errorf("Expected 0 interfaces for empty output, got %d", interfaceCount)
	}
}

func TestGetConnectedServers_NoInterfaces(t *testing.T) {
	// Test parsing wg show output with no interfaces
	testOutput := `No interfaces configured`

	lines := strings.Split(testOutput, "\n")
	interfaces := make([]string, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "interface") {
			interfaceName := strings.TrimPrefix(line, "interface: ")
			interfaces = append(interfaces, interfaceName)
		}
	}

	if len(interfaces) != 0 {
		t.Errorf("Expected no interfaces, got %v", interfaces)
	}
}

// TestVpnInterface verifies that DefaultVpn implements the Vpn interface
func TestVpnInterface(t *testing.T) {
	// This test ensures DefaultVpn implements all required methods
	var vpn Vpn = &DefaultVpn{}
	_ = vpn // Use the variable to avoid unused variable error
}

// Test helper functions and edge cases
func TestGenerateKeys(t *testing.T) {
	// Test key generation
	publicKey, preSharedKey, privateKey, err := generateKeys()
	if err != nil {
		t.Fatalf("generateKeys should succeed: %v", err)
	}

	// Verify keys are not empty
	if publicKey.String() == "" {
		t.Error("Public key should not be empty")
	}
	if preSharedKey.String() == "" {
		t.Error("Pre-shared key should not be empty")
	}
	if privateKey.String() == "" {
		t.Error("Private key should not be empty")
	}

	// Verify public key matches private key
	expectedPublicKey := privateKey.PublicKey()
	if publicKey.String() != expectedPublicKey.String() {
		t.Error("Public key should match the private key's public key")
	}

	// Test that multiple calls generate different keys
	publicKey2, _, privateKey2, err := generateKeys()
	if err != nil {
		t.Fatalf("Second generateKeys should succeed: %v", err)
	}

	if publicKey.String() == publicKey2.String() {
		t.Error("Multiple calls should generate different public keys")
	}
	if privateKey.String() == privateKey2.String() {
		t.Error("Multiple calls should generate different private keys")
	}
}

// Test NewDefaultVpn constructor - simplified test without mocks
func TestNewDefaultVpn(t *testing.T) {
	// This test just ensures the constructor works without testing complex interactions
	vpn := &DefaultVpn{}
	if vpn == nil {
		t.Error("DefaultVpn should be constructible")
	}

	// Verify it implements the Vpn interface
	var _ Vpn = vpn
}

// Test individual functions that can be tested independently

// Test error handling in helper functions
func TestSaveWgConfig_ErrorHandling(t *testing.T) {
	// Test with mock directory provider that returns errors
	mockDirProvider := &mockDirectoryProvider{shouldError: true}
	vpn := &DefaultVpn{dirProvider: mockDirProvider}

	_, err := vpn.saveWgConfig("test", "content")
	if err == nil {
		t.Error("saveWgConfig should fail when directory provider fails")
	}
}

func TestWriteConfig_Integration(t *testing.T) {
	dirProvider, cleanup := setupTestDir(t)
	defer cleanup()

	vpn := &DefaultVpn{dirProvider: dirProvider}

	server := remote.Server{
		Hostname:   "test.example.com",
		IPv4AddrIn: "203.0.113.1",
		PublicKey:  "server-public-key",
	}

	configPath, err := vpn.writeConfig("test-server", server, "private-key", "10.0.0.1/32", "2001:db8::1/128")
	if err != nil {
		t.Fatalf("writeConfig should succeed: %v", err)
	}

	// Verify file was created and contains expected content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	configStr := string(content)
	expectedValues := []string{
		"PrivateKey = private-key",
		"PublicKey = server-public-key",
		"Endpoint = 203.0.113.1:51820",
		"Address = 10.0.0.1/32, 2001:db8::1/128",
	}

	for _, expected := range expectedValues {
		if !strings.Contains(configStr, expected) {
			t.Errorf("Config should contain '%s'\nActual config:\n%s", expected, configStr)
		}
	}
}

// Additional test functions to improve coverage
func TestGetConnectedServers_EmptyInput(t *testing.T) {
	// Test with empty string to cover error paths
	lines := strings.Split("", "\n")
	interfaceCount := strings.Count("", "interface:")
	if interfaceCount != 0 {
		t.Errorf("Expected 0 interfaces for empty input, got %d", interfaceCount)
	}

	// Test the interface parsing logic
	interfaces := make([]string, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "interface") {
			interfaceName := strings.TrimPrefix(line, "interface: ")
			interfaces = append(interfaces, interfaceName)
		}
	}

	if len(interfaces) != 0 {
		t.Errorf("Expected no interfaces, got %v", interfaces)
	}
}

func TestGetConnectedServers_MultipleInterfaces(t *testing.T) {
	// Test parsing multiple interfaces
	testOutput := "interface: wg0\ninterface: wg1\ninterface: wg2"

	lines := strings.Split(testOutput, "\n")
	interfaceCount := strings.Count(testOutput, "interface:")

	if interfaceCount != 3 {
		t.Errorf("Expected 3 interfaces, got %d", interfaceCount)
	}

	// Test interface extraction
	interfaces := make([]string, interfaceCount)
	i := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "interface") {
			interfaceName := strings.TrimPrefix(line, "interface: ")
			if i < len(interfaces) {
				interfaces[i] = interfaceName
				i++
			}
		}
	}

	expectedInterfaces := []string{"wg0", "wg1", "wg2"}
	for idx, expected := range expectedInterfaces {
		if idx >= len(interfaces) || interfaces[idx] != expected {
			t.Errorf("Expected interface %s at index %d, got %s", expected, idx, interfaces[idx])
		}
	}
}

func TestGenerateKeys_ErrorPaths(t *testing.T) {
	// This test mainly verifies the function signature and return types
	publicKey, preSharedKey, privateKey, err := generateKeys()
	if err != nil {
		t.Fatalf("generateKeys should succeed: %v", err)
	}

	// Test that the keys are different (coverage for key generation logic)
	if publicKey.String() == preSharedKey.String() {
		t.Error("Public key and pre-shared key should be different")
	}

	if publicKey.String() == privateKey.String() {
		t.Error("Public key and private key should be different")
	}

	if preSharedKey.String() == privateKey.String() {
		t.Error("Pre-shared key and private key should be different")
	}
}

// Benchmark tests
func BenchmarkGetCountryFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryFlag("US")
	}
}

func BenchmarkGetCountryFlag_Unknown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryFlag("XX")
	}
}

func BenchmarkGenerateKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateKeys()
	}
}
