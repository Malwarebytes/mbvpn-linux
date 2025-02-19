package remote

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Malwarebytes/mbvpn/config"
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

func TestRegisterDevice(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"registerDevice": {"device": {"installationToken": "mock-token"}}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

	token, err := api.RegisterDevice()
	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

func TestActivateDevice(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"activateDevice": {"deviceModules": {"privacy": {"status": "licensed", "termEndsOn": "2025-12-31"}}}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

	deviceModule, err := api.ActivateDevice("mock-token", "mock-key", true)
	assert.NoError(t, err)
	assert.Equal(t, DeviceStatusLicensed, deviceModule.Status)
}

func TestDeactivateDevice(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"deactivateDevice": {"deviceModules": {"privacy": {"status": "unknown", "termEndsOn": "2025-12-31"}}}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

	deviceModule, err := api.DeactivateDevice("mock-token")
	assert.NoError(t, err)
	assert.Equal(t, DeviceStatusUnknown, deviceModule.Status)
}

func TestVpnRegisterPublicKey(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"vpnRegisterPublicKey": {"ipv4": "mock-ipv4", "ipv6": "mock-ipv6", "keyExpirationHours": 24}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

	vpnIpAddresses, err := api.VpnRegisterPublicKey("mock-token", "mock-key")
	assert.NoError(t, err)
	assert.Equal(t, "mock-ipv4", vpnIpAddresses.IpV4)
	assert.Equal(t, "mock-ipv6", vpnIpAddresses.IpV6)
	assert.Equal(t, 24, vpnIpAddresses.KeyExpirationHours)
}

func TestGetVpnNetworkDetails(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"vpnNetworkDetails": {"ip": "mock-ip", "vpnEnabled": true, "geo": {"country": "mock-country", "city": "mock-city", "latitude": 1.23, "longitude": 4.56}}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

	vpnNetworkDetails, err := api.GetVpnNetworkDetails()
	assert.NoError(t, err)
	assert.Equal(t, "mock-ip", vpnNetworkDetails.Ip)
	assert.True(t, vpnNetworkDetails.VpnEnabled)
	assert.Equal(t, "mock-country", vpnNetworkDetails.Geo.Country)
	assert.Equal(t, "mock-city", vpnNetworkDetails.Geo.City)
	assert.Equal(t, float32(1.23), vpnNetworkDetails.Geo.Latitude)
	assert.Equal(t, float32(4.56), vpnNetworkDetails.Geo.Longitude)
}

func TestGetVpnLocations(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"vpnLocations": {"countries": [{"code": "mock-code", "name": "mock-name", "defaultCity": "mock-city", "cities": [{"code": "mock-city-code", "unLocode": "mock-unlocode", "name": "mock-city-name", "latitude": 1.23, "longitude": 4.56, "servers": [{"hostname": "mock-hostname", "ipv4AddrIn": "mock-ipv4", "ipv6AddrIn": "mock-ipv6", "weight": 1, "publicKey": "mock-public-key", "portRanges": [{"from": 1000, "to": 2000}], "multihopPort": 3000, "ipv4Gateway": "mock-ipv4-gateway", "ipv6Gateway": "mock-ipv6-gateway"}]}]}]}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

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

func TestGetVpnClientDefaults(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"vpnClientDefaults": {"ipv4": {"maskAll": "true", "mtu": 1500, "maskNolocal": "false"}, "ipv6": {"maskAll": "true", "mtu": 1500, "maskNolocal": "false"}, "keyExpirationHours": 24}}}`))
	}))
	defer mockServer.Close()

	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

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

// func TestDoRequest(t *testing.T) {
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"data": {"response": "mock-response"}}`))
// 	}))
// 	defer mockServer.Close()

// 	api := NewMockHolocron(mockServer.URL, &MockMachineIdProvider{})

// 	req, err := http.NewRequest("GET", mockServer.URL+"/mock-endpoint", nil)
// 	assert.NoError(t, err)

// 	var result map[string]interface{}
// 	err = api.doRequest("", req, &result)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "mock-response", result["data"].(map[string]interface{})["response"])
// }
