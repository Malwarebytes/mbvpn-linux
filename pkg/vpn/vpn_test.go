package vpn

import (
	"os"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
)

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

// Test NewDefaultVpn constructor - requires root/CAP_NET_ADMIN
func TestNewDefaultVpn(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root or CAP_NET_ADMIN")
	}

	dirProvider, err := config.NewDefaultDirectoryProvider()
	if err != nil {
		t.Fatalf("Failed to create directory provider: %v", err)
	}

	_, err = NewDefaultVpn(nil, nil, nil, dirProvider)
	if err != nil {
		t.Fatalf("NewDefaultVpn should succeed with root/CAP_NET_ADMIN: %v", err)
	}
}

func TestSanitizeInterfaceName(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		expected string
	}{
		{"simple hostname", "server1.example.com", "server1"},
		{"hostname with dashes", "us-east-01.vpn.example.com", "us-east-01"},
		{"long hostname", "very-long-server-name-that-exceeds-limit.example.com", "very-long-serve"},
		{"special chars", "server!@#$%^&*().example.com", "server"},
		{"underscores", "server_01.example.com", "server_01"},
		{"mixed case", "ServerName.example.com", "ServerName"},
		{"numbers only", "12345.example.com", "12345"},
		{"empty after sanitize", "!@#$.example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeInterfaceName(tt.hostname)
			if result != tt.expected {
				t.Errorf("sanitizeInterfaceName(%s) = %s, expected %s", tt.hostname, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"zero bytes", 0, "0 B"},
		{"small bytes", 512, "512 B"},
		{"one KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"one MB", 1024 * 1024, "1.0 MB"},
		{"one GB", 1024 * 1024 * 1024, "1.0 GB"},
		{"one TB", 1024 * 1024 * 1024 * 1024, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %s, expected %s", tt.bytes, result, tt.expected)
			}
		})
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

func BenchmarkSanitizeInterfaceName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sanitizeInterfaceName("very-long-server-name.example.com")
	}
}
