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

func TestRegisterDevice_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"registerDevice": {"device": {"installationToken": "mock-token"}}}}`)
	defer mockServer.Close()

	token, err := api.RegisterDevice()
	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestRegisterDevice_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.RegisterDevice()
	assert.Error(t, err)
}

func TestActivateDevice_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"activateDevice": {"deviceModules": {"privacy": {"status": "licensed", "termEndsOn": "2025-12-31"}}}}}`)
	defer mockServer.Close()

	deviceModule, err := api.ActivateDevice("mock-token", "mock-key", true)
	assert.NoError(t, err)
	assert.Equal(t, DeviceStatusLicensed, deviceModule.Status)
}

func TestActivateDevice_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.ActivateDevice("", "", false)
	assert.Error(t, err)
}

func TestDeactivateDevice_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"deactivateDevice": {"deviceModules": {"privacy": {"status": "unknown", "termEndsOn": "2025-12-31"}}}}}`)
	defer mockServer.Close()

	deviceModule, err := api.DeactivateDevice("mock-token")
	assert.NoError(t, err)
	assert.Equal(t, DeviceStatusUnknown, deviceModule.Status)
}

func TestDeactivateDevice_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.DeactivateDevice("")
	assert.Error(t, err)
}

func TestVpnRegisterPublicKey_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"vpnRegisterPublicKey": {"ipv4": "mock-ipv4", "ipv6": "mock-ipv6", "keyExpirationHours": 24}}}`)
	defer mockServer.Close()

	vpnIpAddresses, err := api.VpnRegisterPublicKey("mock-token", "mock-key")
	assert.NoError(t, err)
	assert.Equal(t, "mock-ipv4", vpnIpAddresses.IpV4)
	assert.Equal(t, "mock-ipv6", vpnIpAddresses.IpV6)
	assert.Equal(t, 24, vpnIpAddresses.KeyExpirationHours)
}

func TestVpnRegisterPublicKey_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.VpnRegisterPublicKey("", "")
	assert.Error(t, err)
}

func TestGetVpnNetworkDetails_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"vpnNetworkDetails": {"ip": "mock-ip", "vpnEnabled": true, "geo": {"country": "mock-country", "city": "mock-city", "latitude": 1.23, "longitude": 4.56}}}}`)
	defer mockServer.Close()

	vpnNetworkDetails, err := api.GetVpnNetworkDetails()
	assert.NoError(t, err)
	assert.Equal(t, "mock-ip", vpnNetworkDetails.Ip)
	assert.True(t, vpnNetworkDetails.VpnEnabled)
	assert.Equal(t, "mock-country", vpnNetworkDetails.Geo.Country)
	assert.Equal(t, "mock-city", vpnNetworkDetails.Geo.City)
	assert.Equal(t, float32(1.23), vpnNetworkDetails.Geo.Latitude)
	assert.Equal(t, float32(4.56), vpnNetworkDetails.Geo.Longitude)
}

func TestGetVpnNetworkDetails_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.GetVpnNetworkDetails()
	assert.Error(t, err)
}

func TestGetVpnLocations_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"vpnLocations": {"countries": [{"code": "mock-code", "name": "mock-name", "defaultCity": "mock-city", "cities": [{"code": "mock-city-code", "unLocode": "mock-unlocode", "name": "mock-city-name", "latitude": 1.23, "longitude": 4.56, "servers": [{"hostname": "mock-hostname", "ipv4AddrIn": "mock-ipv4", "ipv6AddrIn": "mock-ipv6", "weight": 1, "publicKey": "mock-public-key", "portRanges": [{"from": 1000, "to": 2000}], "multihopPort": 3000, "ipv4Gateway": "mock-ipv4-gateway", "ipv6Gateway": "mock-ipv6-gateway"}]}]}]}}}`)
	defer mockServer.Close()

	vpnLocations, err := api.GetVpnLocations()
	assert.NoError(t, err)
	assert.Equal(t, "mock-code", vpnLocations.Countries[0].Code)
	assert.Equal(t, "mock-name", vpnLocations.Countries[0].Name)
	assert.Equal(t, "mock-city", vpnLocations.Countries[0].DefaultCity)
	assert.Equal(t, "mock-city-code", vpnLocations.Countries[0].Cities[0].Code)
	assert.Equal(t, "mock-unlocode", vpnLocations.Countries[0].Cities[0].UnLocode)
	assert.Equal(t, "mock-city-name", vpnLocations.Countries[0].Cities[0].Name)
	assert.Equal(t, float32(1.23), vpnLocations.Countries[0].Cities[0].Latitude)
	assert.Equal(t, float32(4.56), vpnLocations.Countries[0].Cities[0].Longitude)
	assert.Equal(t, "mock-hostname", vpnLocations.Countries[0].Cities[0].Servers[0].Hostname)
	assert.Equal(t, "mock-ipv4", vpnLocations.Countries[0].Cities[0].Servers[0].IPv4AddrIn)
	assert.Equal(t, "mock-ipv6", vpnLocations.Countries[0].Cities[0].Servers[0].IPv6AddrIn)
	assert.Equal(t, 1, vpnLocations.Countries[0].Cities[0].Servers[0].Weight)
	assert.Equal(t, "mock-public-key", vpnLocations.Countries[0].Cities[0].Servers[0].PublicKey)
	assert.Equal(t, 1000, vpnLocations.Countries[0].Cities[0].Servers[0].PortRanges[0].From)
	assert.Equal(t, 2000, vpnLocations.Countries[0].Cities[0].Servers[0].PortRanges[0].To)
	assert.Equal(t, 3000, vpnLocations.Countries[0].Cities[0].Servers[0].MultihopPort)
	assert.Equal(t, "mock-ipv4-gateway", vpnLocations.Countries[0].Cities[0].Servers[0].IPv4Gateway)
	assert.Equal(t, "mock-ipv6-gateway", vpnLocations.Countries[0].Cities[0].Servers[0].IPv6Gateway)
}

func TestGetVpnLocations_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.GetVpnLocations()
	assert.Error(t, err)
}

func TestGetVpnClientDefaults_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"vpnClientDefaults": {"ipv4": {"maskAll": "true", "mtu": 1500, "maskNolocal": "false"}, "ipv6": {"maskAll": "true", "mtu": 1500, "maskNolocal": "false"}, "keyExpirationHours": 24}}}`)
	defer mockServer.Close()

	vpnClientDefaults, err := api.GetVpnClientDefaults()
	assert.NoError(t, err)
	assert.Equal(t, "true", vpnClientDefaults.IPv4.MaskAll)
	assert.Equal(t, 1500, vpnClientDefaults.IPv4.MTU)
	assert.Equal(t, "false", vpnClientDefaults.IPv4.MaskNolocal)
	assert.Equal(t, "true", vpnClientDefaults.IPv6.MaskAll)
	assert.Equal(t, 1500, vpnClientDefaults.IPv6.MTU)
	assert.Equal(t, "false", vpnClientDefaults.IPv6.MaskNolocal)
	assert.Equal(t, 24, vpnClientDefaults.KeyExpirationHours)
}

func TestGetVpnClientDefaults_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	_, err := api.GetVpnClientDefaults()
	assert.Error(t, err)
}

func TestDoRequest_Success(t *testing.T) {
	mockServer, api := setupMockServer(`{"data": {"response": "mock-response"}}`)
	defer mockServer.Close()

	var body map[string]interface{}
	var result map[string]interface{}
	err := api.doRequest("", &body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "mock-response", result["data"].(map[string]interface{})["response"])
}

func TestDoRequest_Error(t *testing.T) {
	mockServer, api := setupMockServer(`{ "errors": [ { "message": "mock-error" } ] }`)
	defer mockServer.Close()

	var body map[string]interface{}
	var result map[string]interface{}
	err := api.doRequest("", &body, &result)
	assert.Error(t, err)
}

func TestDoRequest_InvalidJson(t *testing.T) {
	mockServer, api := setupMockServer(`{`)
	defer mockServer.Close()

	var body map[string]interface{}
	var result map[string]interface{}
	err := api.doRequest("", &body, &result)
	assert.Error(t, err)
}
