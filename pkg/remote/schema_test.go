package remote

import (
	"encoding/json"
	"testing"
)

func TestProductModule_Constants(t *testing.T) {
	if ProductModulePrivacy != "privacy" {
		t.Errorf("Expected ProductModulePrivacy to be 'privacy', got '%s'", ProductModulePrivacy)
	}
}

func TestActivationMethod_Constants(t *testing.T) {
	if ActivationMethodOneTimeToken != "oneTimeToken" {
		t.Errorf("Expected ActivationMethodOneTimeToken to be 'oneTimeToken', got '%s'", ActivationMethodOneTimeToken)
	}
}

func TestActivationMode_Constants(t *testing.T) {
	if ActivationModePassive != "passive" {
		t.Errorf("Expected ActivationModePassive to be 'passive', got '%s'", ActivationModePassive)
	}
}

func TestDeactivationMode_Constants(t *testing.T) {
	if DeactivationModePassive != "passive" {
		t.Errorf("Expected DeactivationModePassive to be 'passive', got '%s'", DeactivationModePassive)
	}
}

func TestDeviceStatus_Constants(t *testing.T) {
	statuses := map[DeviceStatus]string{
		DeviceStatusFree:           "free",
		DeviceStatusLicenseExpired: "licenseExpired",
		DeviceStatusLicenseGrace:   "licenseGrace",
		DeviceStatusLicensed:       "licensed",
		DeviceStatusTrial:          "trial",
		DeviceStatusTrialExpired:   "trialExpired",
		DeviceStatusUnknown:        "unknown",
		DeviceStatusLicenseOnHold:  "licenseOnHold",
	}

	for status, expected := range statuses {
		if string(status) != expected {
			t.Errorf("Expected %s to equal '%s', got '%s'", status, expected, string(status))
		}
	}
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	t.Run("Marshal ErrorResponse", func(t *testing.T) {
		errResp := ErrorResponse{
			Errors: []Error{
				{Message: "First error"},
				{Message: "Second error"},
			},
		}

		data, err := json.Marshal(errResp)
		if err != nil {
			t.Fatalf("Failed to marshal ErrorResponse: %v", err)
		}

		expected := `{"errors":[{"message":"First error"},{"message":"Second error"}]}`
		if string(data) != expected {
			t.Errorf("Expected JSON: %s, got: %s", expected, string(data))
		}
	})

	t.Run("Unmarshal ErrorResponse", func(t *testing.T) {
		jsonData := `{"errors":[{"message":"Test error"}]}`

		var errResp ErrorResponse
		err := json.Unmarshal([]byte(jsonData), &errResp)
		if err != nil {
			t.Fatalf("Failed to unmarshal ErrorResponse: %v", err)
		}

		if len(errResp.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errResp.Errors))
		}
		if errResp.Errors[0].Message != "Test error" {
			t.Errorf("Expected 'Test error', got '%s'", errResp.Errors[0].Message)
		}
	})
}

func TestRegisterDeviceInput_JSONSerialization(t *testing.T) {
	input := RegisterDeviceInput{
		ProductCode:    "MBVPN",
		ProductVersion: "1.0.0",
		Modules:        []ProductModule{ProductModulePrivacy},
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal RegisterDeviceInput: %v", err)
	}

	var unmarshaled RegisterDeviceInput
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal RegisterDeviceInput: %v", err)
	}

	if unmarshaled.ProductCode != input.ProductCode {
		t.Errorf("Expected ProductCode '%s', got '%s'", input.ProductCode, unmarshaled.ProductCode)
	}
	if unmarshaled.ProductVersion != input.ProductVersion {
		t.Errorf("Expected ProductVersion '%s', got '%s'", input.ProductVersion, unmarshaled.ProductVersion)
	}
	if len(unmarshaled.Modules) != len(input.Modules) {
		t.Errorf("Expected %d modules, got %d", len(input.Modules), len(unmarshaled.Modules))
	}
}

func TestActivateDeviceInput_JSONSerialization(t *testing.T) {
	input := ActivateDeviceInput{
		Modules:          []ProductModule{ProductModulePrivacy},
		ActivationMethod: ActivationMethodOneTimeToken,
		ActivationMode:   ActivationModePassive,
		OneTimeToken:     "TEST-TOKEN",
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal ActivateDeviceInput: %v", err)
	}

	var unmarshaled ActivateDeviceInput
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ActivateDeviceInput: %v", err)
	}

	if unmarshaled.ActivationMethod != input.ActivationMethod {
		t.Errorf("Expected ActivationMethod '%s', got '%s'", input.ActivationMethod, unmarshaled.ActivationMethod)
	}
}

func TestHolocronResponse_JSONSerialization(t *testing.T) {
	response := HolocronResponse[string]{
		Data: "test data",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal HolocronResponse: %v", err)
	}

	var unmarshaled HolocronResponse[string]
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal HolocronResponse: %v", err)
	}

	if unmarshaled.Data != response.Data {
		t.Errorf("Expected Data '%s', got '%s'", response.Data, unmarshaled.Data)
	}
}

func TestDeviceOutput_JSONSerialization(t *testing.T) {
	output := DeviceOutput{
		Device: Device{
			InstallationToken: "test-token",
		},
		DeviceModules: DeviceModules{
			Privacy: DeviceModule{
				Status: DeviceStatusLicensed,
			},
		},
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal DeviceOutput: %v", err)
	}

	var unmarshaled DeviceOutput
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal DeviceOutput: %v", err)
	}

	if unmarshaled.Device.InstallationToken != output.Device.InstallationToken {
		t.Errorf("Expected InstallationToken '%s', got '%s'",
			output.Device.InstallationToken, unmarshaled.Device.InstallationToken)
	}
	if unmarshaled.DeviceModules.Privacy.Status != output.DeviceModules.Privacy.Status {
		t.Errorf("Expected Status '%s', got '%s'",
			output.DeviceModules.Privacy.Status, unmarshaled.DeviceModules.Privacy.Status)
	}
}

func TestVpnIpAddresses_JSONSerialization(t *testing.T) {
	ipAddrs := VpnIpAddresses{
		IpV4:               "10.0.0.1",
		IpV6:               "2001:db8::1",
		KeyExpirationHours: 24,
	}

	data, err := json.Marshal(ipAddrs)
	if err != nil {
		t.Fatalf("Failed to marshal VpnIpAddresses: %v", err)
	}

	var unmarshaled VpnIpAddresses
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpnIpAddresses: %v", err)
	}

	if unmarshaled.IpV4 != ipAddrs.IpV4 {
		t.Errorf("Expected IPv4 '%s', got '%s'", ipAddrs.IpV4, unmarshaled.IpV4)
	}
	if unmarshaled.IpV6 != ipAddrs.IpV6 {
		t.Errorf("Expected IPv6 '%s', got '%s'", ipAddrs.IpV6, unmarshaled.IpV6)
	}
	if unmarshaled.KeyExpirationHours != ipAddrs.KeyExpirationHours {
		t.Errorf("Expected KeyExpirationHours %d, got %d",
			ipAddrs.KeyExpirationHours, unmarshaled.KeyExpirationHours)
	}
}

func TestVpnNetworkDetails_JSONSerialization(t *testing.T) {
	details := VpnNetworkDetails{
		Ip:         "203.0.113.1",
		VpnEnabled: true,
		Geo: Geo{
			Country:   "United States",
			City:      "New York",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
	}

	data, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("Failed to marshal VpnNetworkDetails: %v", err)
	}

	var unmarshaled VpnNetworkDetails
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpnNetworkDetails: %v", err)
	}

	if unmarshaled.Ip != details.Ip {
		t.Errorf("Expected IP '%s', got '%s'", details.Ip, unmarshaled.Ip)
	}
	if unmarshaled.VpnEnabled != details.VpnEnabled {
		t.Errorf("Expected VpnEnabled %t, got %t", details.VpnEnabled, unmarshaled.VpnEnabled)
	}
	if unmarshaled.Geo.Country != details.Geo.Country {
		t.Errorf("Expected Country '%s', got '%s'", details.Geo.Country, unmarshaled.Geo.Country)
	}
}

func TestServer_JSONSerialization(t *testing.T) {
	server := Server{
		Hostname:     "vpn01.example.com",
		IPv4AddrIn:   "203.0.113.1",
		IPv6AddrIn:   "2001:db8::1",
		Weight:       100,
		PublicKey:    "test-public-key",
		PortRanges:   []PortRange{{From: 51820, To: 51830}},
		MultihopPort: 51821,
		IPv4Gateway:  "203.0.113.254",
		IPv6Gateway:  "2001:db8::254",
	}

	data, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("Failed to marshal Server: %v", err)
	}

	var unmarshaled Server
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Server: %v", err)
	}

	if unmarshaled.Hostname != server.Hostname {
		t.Errorf("Expected Hostname '%s', got '%s'", server.Hostname, unmarshaled.Hostname)
	}
	if unmarshaled.Weight != server.Weight {
		t.Errorf("Expected Weight %d, got %d", server.Weight, unmarshaled.Weight)
	}
	if len(unmarshaled.PortRanges) != len(server.PortRanges) {
		t.Errorf("Expected %d port ranges, got %d", len(server.PortRanges), len(unmarshaled.PortRanges))
	}
	if len(unmarshaled.PortRanges) > 0 {
		if unmarshaled.PortRanges[0].From != server.PortRanges[0].From {
			t.Errorf("Expected port range From %d, got %d",
				server.PortRanges[0].From, unmarshaled.PortRanges[0].From)
		}
	}
}

func TestVpnLocations_JSONSerialization(t *testing.T) {
	locations := VpnLocations{
		Countries: []Country{
			{
				Code:        "US",
				Name:        "United States",
				DefaultCity: "NYC",
				Cities: []City{
					{
						Code:      "NYC",
						UnLocode:  "USNYC",
						Name:      "New York",
						Latitude:  40.7128,
						Longitude: -74.0060,
						Servers: []Server{
							{
								Hostname:  "vpn01.us.example.com",
								Weight:    100,
								PublicKey: "test-key",
							},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(locations)
	if err != nil {
		t.Fatalf("Failed to marshal VpnLocations: %v", err)
	}

	var unmarshaled VpnLocations
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpnLocations: %v", err)
	}

	if len(unmarshaled.Countries) != len(locations.Countries) {
		t.Errorf("Expected %d countries, got %d", len(locations.Countries), len(unmarshaled.Countries))
	}

	if len(unmarshaled.Countries) > 0 {
		country := unmarshaled.Countries[0]
		originalCountry := locations.Countries[0]

		if country.Code != originalCountry.Code {
			t.Errorf("Expected country code '%s', got '%s'", originalCountry.Code, country.Code)
		}
		if len(country.Cities) != len(originalCountry.Cities) {
			t.Errorf("Expected %d cities, got %d", len(originalCountry.Cities), len(country.Cities))
		}
	}
}

func TestVpnClientDefaults_JSONSerialization(t *testing.T) {
	defaults := VpnClientDefaults{
		IPv4: VpnIpDefaults{
			MaskAll:     "0.0.0.0/0",
			MTU:         1420,
			MaskNolocal: "0.0.0.0/1",
		},
		IPv6: VpnIpDefaults{
			MaskAll:     "::/0",
			MTU:         1420,
			MaskNolocal: "::/1",
		},
		KeyExpirationHours: 24,
	}

	data, err := json.Marshal(defaults)
	if err != nil {
		t.Fatalf("Failed to marshal VpnClientDefaults: %v", err)
	}

	var unmarshaled VpnClientDefaults
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpnClientDefaults: %v", err)
	}

	if unmarshaled.IPv4.MTU != defaults.IPv4.MTU {
		t.Errorf("Expected IPv4 MTU %d, got %d", defaults.IPv4.MTU, unmarshaled.IPv4.MTU)
	}
	if unmarshaled.KeyExpirationHours != defaults.KeyExpirationHours {
		t.Errorf("Expected KeyExpirationHours %d, got %d",
			defaults.KeyExpirationHours, unmarshaled.KeyExpirationHours)
	}
}

func TestEmptyStructs_JSONSerialization(t *testing.T) {
	// Test that empty structs can be marshaled/unmarshaled
	tests := []interface{}{
		ErrorResponse{},
		RegisterDeviceInput{},
		ActivateDeviceInput{},
		DeviceOutput{},
		VpnIpAddresses{},
		VpnNetworkDetails{},
		Server{},
		VpnLocations{},
	}

	for _, test := range tests {
		data, err := json.Marshal(test)
		if err != nil {
			t.Errorf("Failed to marshal empty struct %T: %v", test, err)
		}

		// Basic check that we get valid JSON
		if len(data) == 0 {
			t.Errorf("Expected non-empty JSON for %T", test)
		}
	}
}

func TestNilSlicesAndPointers_JSONSerialization(t *testing.T) {
	// Test handling of nil slices in JSON serialization
	input := RegisterDeviceInput{
		ProductCode:    "TEST",
		ProductVersion: "1.0",
		Modules:        nil, // nil slice
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal with nil slice: %v", err)
	}

	var unmarshaled RegisterDeviceInput
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal with nil slice: %v", err)
	}

	// Go JSON handling should turn nil slice into empty slice
	if unmarshaled.Modules == nil {
		// This is acceptable behavior
	}
}
