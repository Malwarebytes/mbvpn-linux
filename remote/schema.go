package remote

type ProductModule string

const (
	ProductModulePrivacy ProductModule = "privacy"
)

type ActivationMethod string

const (
	ActivationMethodOneTimeToken ActivationMethod = "oneTimeToken"
	ActivationMethodLicenseKey   ActivationMethod = "licenseKey"
)

type ActivationMode string

const (
	ActivationModePassive ActivationMode = "passive"
)

type DeactivationMode string

const (
	DeactivationModePassive DeactivationMode = "passive"
)

type RegisterDeviceInput struct {
	ProductCode    string          `json:"productCode"`
	ProductVersion string          `json:"productVersion"`
	Modules        []ProductModule `json:"modules"`
}

type ActivateDeviceInput struct {
	Modules          []ProductModule  `json:"modules"`
	ActivationMethod ActivationMethod `json:"activationMethod"`
	ActivationMode   ActivationMode   `json:"activationMode"`
	OneTimeToken     string           `json:"oneTimeToken"`
	LicenseKey       string           `json:"licenseKey"`
}

type ActivateDeviceResponse struct {
	Data ActivateDeviceResponseData `json:"data"`
}

type ActivateDeviceResponseData struct {
	ActivateDevice DeviceOutput `json:"activateDevice"`
}

type DeactivateDeviceInput struct {
	Modules          []ProductModule  `json:"modules"`
	DeactivationMode DeactivationMode `json:"deactivationMode"`
}

type DeactivateDeviceResponse struct {
	Data DeactivateDeviceResponseData `json:"data"`
}

type DeactivateDeviceResponseData struct {
	DeactivateDevice DeviceOutput `json:"deactivateDevice"`
}

type RegisterDeviceResponse struct {
	Data RegisterDeviceResponseData `json:"data"`
}

type RegisterDeviceResponseData struct {
	RegisterDevice DeviceOutput `json:"registerDevice"`
}

type DeviceOutput struct {
	Device        Device        `json:"device"`
	DeviceModules DeviceModules `json:"deviceModules"`
}

type Device struct {
	InstallationToken string `json:"installationToken"`
}

type DeviceModules struct {
	Privacy DeviceModule `json:"privacy"`
}

type DeviceModule struct {
	Status DeviceStatus `json:"status"`
}

type DeviceStatus string

const (
	DeviceStatusFree           DeviceStatus = "free"
	DeviceStatusLicenseExpired DeviceStatus = "licenseExpired"
	DeviceStatusLicenseGrace   DeviceStatus = "licenseGrace"
	DeviceStatusLicensed       DeviceStatus = "licensed"
	DeviceStatusTrial          DeviceStatus = "trial"
	DeviceStatusTrialExpired   DeviceStatus = "trialExpired"
	DeviceStatusUnknown        DeviceStatus = "unknown"
	DeviceStatusLicenseOnHold  DeviceStatus = "licenseOnHold"
)

type VpnRegisterPublicKeyInput struct {
	PublicKey string `json:"publicKey"`
}

type VpnRegisterPublicKeyResponse struct {
	Data VpnRegisterPublicKeyResponseData `json:"data"`
}

type VpnRegisterPublicKeyResponseData struct {
	VpnIpAddresses VpnIpAddresses `json:"vpnRegisterPublicKey"`
}

type VpnIpAddresses struct {
	IpV4               string `json:"ipv4"`
	IpV6               string `json:"ipv6"`
	KeyExpirationHours int    `json:"keyExpirationHours"`
}

type VpnNetworkDetailsResponse struct {
	Data VpnNetworkDetailsResponseData `json:"data"`
}

type VpnNetworkDetailsResponseData struct {
	Details VpnNetworkDetails `json:"vpnNetworkDetails"`
}

type VpnNetworkDetails struct {
	Ip         string `json:"ip"`
	VpnEnabled bool   `json:"vpnEnabled"`
	Geo        Geo    `json:"geo"`
}

type Geo struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type VpnLocationsResponse struct {
	Data VpnLocationsResponseData `json:"data"`
}

type VpnLocationsResponseData struct {
	VpnLocations VpnLocations `json:"vpnLocations"`
}

type VpnLocations struct {
	Countries []Country `json:"countries"`
}

type Country struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	DefaultCity string `json:"defaultCity"`
	Cities      []City `json:"cities"`
}

type City struct {
	Code      string   `json:"code"`
	UnLocode  string   `json:"unLocode"`
	Name      string   `json:"name"`
	Latitude  float32  `json:"latitude"`
	Longitude float32  `json:"longitude"`
	Servers   []Server `json:"servers"`
}

type Server struct {
	Hostname     string      `json:"hostname"`
	IPv4AddrIn   string      `json:"ipv4AddrIn"`
	IPv6AddrIn   string      `json:"ipv6AddrIn"`
	Weight       int         `json:"weight"`
	PublicKey    string      `json:"publicKey"`
	PortRanges   []PortRange `json:"portRanges"`
	MultihopPort int         `json:"multihopPort"`
	IPv4Gateway  string      `json:"ipv4Gateway"`
	IPv6Gateway  string      `json:"ipv6Gateway"`
}

type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type VpnClientDefaultsResponse struct {
	Data VpnClientDefaultsResponseData `json:"data"`
}

type VpnClientDefaultsResponseData struct {
	VpnClientDefaults VpnClientDefaults `json:"vpnClientDefaults"`
}

type VpnClientDefaults struct {
	IPv4           VpnIpDefaults `json:"ipv4"`
	IPv6           VpnIpDefaults `json:"ipv6"`
	KeyExpiryHours int           `json:"keyExpirationHours"`
}

type VpnIpDefaults struct {
	MaskAll     string `json:"maskAll"`
	MTU         int    `json:"mtu"`
	MaskNolocal string `json:"maskNolocal"`
}
