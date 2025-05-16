package remote

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/stretchr/testify/assert"
)

type MockMachineIdProvider struct{}

func (m *MockMachineIdProvider) Get() (string, error) {
	return "mock-machine-id", nil
}

func NewMockHolocron(mockUrl string, mIdProvider config.MachineIdProvider) *DefaultHolocron {
	return &DefaultHolocron{
		baseUrl:     mockUrl,
		mIdProvider: mIdProvider,
	}
}

func setupMockServer(response string) (*httptest.Server, *DefaultHolocron) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})
	return mockServer, api
}

func TestApiSuccessCases(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		testFunc func(*DefaultHolocron) (interface{}, error)
		validate func(t *testing.T, result interface{})
	}{
		{
			name:     "RegisterDevice",
			response: `{"data": {"registerDevice": {"device": {"installationToken": "mock-token"}}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.RegisterDevice()
			},
			validate: func(t *testing.T, result interface{}) {
				token := result.(string)
				assert.Equal(t, "mock-token", token)
			},
		},
		{
			name:     "ActivateDevice",
			response: `{"data": {"activateDevice": {"deviceModules": {"privacy": {"status": "licensed", "termEndsOn": "2025-12-31"}}}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.ActivateDevice("mock-token", "mock-key", true)
			},
			validate: func(t *testing.T, result interface{}) {
				deviceModule := result.(*DeviceModule)
				assert.Equal(t, DeviceStatusLicensed, deviceModule.Status)
			},
		},
		{
			name:     "DeactivateDevice",
			response: `{"data": {"deactivateDevice": {"deviceModules": {"privacy": {"status": "unknown", "termEndsOn": "2025-12-31"}}}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.DeactivateDevice("mock-token")
			},
			validate: func(t *testing.T, result interface{}) {
				deviceModule := result.(*DeviceModule)
				assert.Equal(t, DeviceStatusUnknown, deviceModule.Status)
			},
		},
		{
			name:     "VpnRegisterPublicKey",
			response: `{"data": {"vpnRegisterPublicKey": {"ipv4": "mock-ipv4", "ipv6": "mock-ipv6", "keyExpirationHours": 24}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.VpnRegisterPublicKey("mock-token", "mock-key")
			},
			validate: func(t *testing.T, result interface{}) {
				vpnIpAddresses := result.(*VpnIpAddresses)
				assert.Equal(t, "mock-ipv4", vpnIpAddresses.IpV4)
				assert.Equal(t, "mock-ipv6", vpnIpAddresses.IpV6)
				assert.Equal(t, 24, vpnIpAddresses.KeyExpirationHours)
			},
		},
		{
			name:     "GetVpnNetworkDetails",
			response: `{"data": {"vpnNetworkDetails": {"ip": "mock-ip", "vpnEnabled": true, "geo": {"country": "mock-country", "city": "mock-city", "latitude": 1.23, "longitude": 4.56}}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.GetVpnNetworkDetails()
			},
			validate: func(t *testing.T, result interface{}) {
				details := result.(*VpnNetworkDetails)
				assert.Equal(t, "mock-ip", details.Ip)
				assert.True(t, details.VpnEnabled)
				assert.Equal(t, "mock-country", details.Geo.Country)
				assert.Equal(t, "mock-city", details.Geo.City)
				assert.Equal(t, float32(1.23), details.Geo.Latitude)
				assert.Equal(t, float32(4.56), details.Geo.Longitude)
			},
		},
		{
			name:     "GetVpnLocations",
			response: `{"data": {"vpnLocations": {"countries": [{"code": "mock-code", "name": "mock-name", "defaultCity": "mock-city", "cities": [{"code": "mock-city-code", "unLocode": "mock-unlocode", "name": "mock-city-name", "latitude": 1.23, "longitude": 4.56, "servers": [{"hostname": "mock-hostname", "ipv4AddrIn": "mock-ipv4", "ipv6AddrIn": "mock-ipv6", "weight": 1, "publicKey": "mock-public-key", "portRanges": [{"from": 1000, "to": 2000}], "multihopPort": 3000, "ipv4Gateway": "mock-ipv4-gateway", "ipv6Gateway": "mock-ipv6-gateway"}]}]}]}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.GetVpnLocations()
			},
			validate: func(t *testing.T, result interface{}) {
				vpnLocations := result.(*VpnLocations)
				assert.Equal(t, "mock-code", vpnLocations.Countries[0].Code)
				assert.Equal(t, "mock-name", vpnLocations.Countries[0].Name)
				assert.Equal(t, "mock-city", vpnLocations.Countries[0].DefaultCity)
				assert.Equal(t, "mock-city-code", vpnLocations.Countries[0].Cities[0].Code)
				assert.Equal(t, "mock-hostname", vpnLocations.Countries[0].Cities[0].Servers[0].Hostname)
			},
		},
		{
			name:     "GetVpnClientDefaults",
			response: `{"data": {"vpnClientDefaults": {"ipv4": {"maskAll": "true", "mtu": 1500, "maskNolocal": "false"}, "ipv6": {"maskAll": "true", "mtu": 1500, "maskNolocal": "false"}, "keyExpirationHours": 24}}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.GetVpnClientDefaults()
			},
			validate: func(t *testing.T, result interface{}) {
				defaults := result.(*VpnClientDefaults)
				assert.Equal(t, "true", defaults.IPv4.MaskAll)
				assert.Equal(t, 1500, defaults.IPv4.MTU)
				assert.Equal(t, "false", defaults.IPv4.MaskNolocal)
				assert.Equal(t, "true", defaults.IPv6.MaskAll)
				assert.Equal(t, 1500, defaults.IPv6.MTU)
				assert.Equal(t, "false", defaults.IPv6.MaskNolocal)
				assert.Equal(t, 24, defaults.KeyExpirationHours)
			},
		},
		{
			name:     "DoRequest",
			response: `{"data": {"response": "mock-response"}}`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				var body map[string]interface{}
				var result map[string]interface{}
				err := api.doRequest("", &body, &result)
				return result, err
			},
			validate: func(t *testing.T, result interface{}) {
				resultMap := result.(map[string]interface{})
				assert.Equal(t, "mock-response", resultMap["data"].(map[string]interface{})["response"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer, api := setupMockServer(tc.response)
			defer mockServer.Close()

			result, err := tc.testFunc(api)
			assert.NoError(t, err)
			tc.validate(t, result)
		})
	}
}

func TestApiErrorCases(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		testFunc func(*DefaultHolocron) (interface{}, error)
	}{
		{
			name:     "RegisterDevice",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.RegisterDevice()
			},
		},
		{
			name:     "ActivateDevice",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.ActivateDevice("", "", false)
			},
		},
		{
			name:     "DeactivateDevice",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.DeactivateDevice("")
			},
		},
		{
			name:     "VpnRegisterPublicKey",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.VpnRegisterPublicKey("", "")
			},
		},
		{
			name:     "GetVpnNetworkDetails",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.GetVpnNetworkDetails()
			},
		},
		{
			name:     "GetVpnLocations",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.GetVpnLocations()
			},
		},
		{
			name:     "GetVpnClientDefaults",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			testFunc: func(api *DefaultHolocron) (interface{}, error) {
				return api.GetVpnClientDefaults()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer, api := setupMockServer(tc.response)
			defer mockServer.Close()

			_, err := tc.testFunc(api)
			assert.Error(t, err)
		})
	}
}

func TestDoRequestVariants(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name:     "Success",
			response: `{"data": {"response": "mock-response"}}`,
			wantErr:  false,
		},
		{
			name:     "Error",
			response: `{ "errors": [ { "message": "mock-error" } ] }`,
			wantErr:  true,
		},
		{
			name:     "InvalidJson",
			response: `{`,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer, api := setupMockServer(tc.response)
			defer mockServer.Close()

			var body map[string]interface{}
			var result map[string]interface{}
			err := api.doRequest("", &body, &result)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "mock-response", result["data"].(map[string]interface{})["response"])
			}
		})
	}
}
