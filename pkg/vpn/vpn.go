package vpn

import (
	"fmt"
	"net"
	"strings"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	"github.com/malwarebytes/mbvpn-linux/pkg/errors"
	"github.com/malwarebytes/mbvpn-linux/pkg/output"
	"github.com/malwarebytes/mbvpn-linux/pkg/remote"
	"github.com/malwarebytes/mbvpn-linux/pkg/servers"
	"github.com/malwarebytes/mbvpn-linux/pkg/wireguard"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Vpn interface {
	Servers(showCities bool, showServers bool) error
	Connect(cfg string) error
	Disconnect(cfg string) error
	Status() error
}

type DefaultVpn struct {
	cfgProvider   config.ConfigProvider
	holocron      remote.Holocron
	serverStorage servers.ServerStorage
	dirProvider   config.DirectoryProvider
	wgManager     wireguard.Manager
}

func NewDefaultVpn(cfgProvider config.ConfigProvider, holocron remote.Holocron, serverStorage servers.ServerStorage, dirProvider config.DirectoryProvider) (Vpn, error) {
	mgr, err := wireguard.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create wireguard manager: %w", err)
	}

	return &DefaultVpn{
		cfgProvider:   cfgProvider,
		holocron:      holocron,
		serverStorage: serverStorage,
		dirProvider:   dirProvider,
		wgManager:     mgr,
	}, nil
}

func (vpn *DefaultVpn) Servers(showCities bool, showServers bool) error {
	locations, err := vpn.holocron.GetVpnLocations()
	if err != nil {
		return errors.NewNetworkError("get server list", err)
	}

	err = vpn.serverStorage.Save(locations)
	if err != nil {
		return errors.NewConfigError("save server list", err)
	}

	for _, country := range locations.Countries {
		countryFlag := getCountryFlag(country.Code)
		output.PrintMsg(fmt.Sprintf("%s %s, %s", countryFlag, country.Name, country.Code), output.MsgOutput)

		for i, city := range country.Cities {
			if showCities {
				isLastCity := i == len(country.Cities)-1
				cityPrefix := "└─"
				if !isLastCity {
					cityPrefix = "├─"
				}
				output.PrintMsg(fmt.Sprintf("  %s %s, %s", cityPrefix, city.Name, city.Code), output.MsgOutput)

				if showServers {
					for j, s := range city.Servers {
						isLastServer := j == len(city.Servers)-1
						serverPrefix := "└─"
						if !isLastServer {
							serverPrefix = "├─"
						}
						output.PrintMsg(fmt.Sprintf("     %s %s", serverPrefix, strings.Split(s.Hostname, ".")[0]), output.MsgOutput)
					}
				}
			}
		}
	}
	return nil
}

// getCountryFlag returns the flag emoji for a given country code
func getCountryFlag(countryCode string) string {
	// Convert country code to uppercase to ensure consistent handling
	countryCode = strings.ToUpper(countryCode)

	// Map of country codes to flag emojis
	countryFlags := map[string]string{
		"AF": "🇦🇫", // Afghanistan
		"AL": "🇦🇱", // Albania
		"DZ": "🇩🇿", // Algeria
		"AR": "🇦🇷", // Argentina
		"AM": "🇦🇲", // Armenia
		"AU": "🇦🇺", // Australia
		"AT": "🇦🇹", // Austria
		"AZ": "🇦🇿", // Azerbaijan
		"BH": "🇧🇭", // Bahrain
		"BD": "🇧🇩", // Bangladesh
		"BY": "🇧🇾", // Belarus
		"BE": "🇧🇪", // Belgium
		"BO": "🇧🇴", // Bolivia
		"BA": "🇧🇦", // Bosnia and Herzegovina
		"BR": "🇧🇷", // Brazil
		"KH": "🇰🇭", // Cambodia
		"CA": "🇨🇦", // Canada
		"CL": "🇨🇱", // Chile
		"CN": "🇨🇳", // China
		"CO": "🇨🇴", // Colombia
		"CR": "🇨🇷", // Costa Rica
		"HR": "🇭🇷", // Croatia
		"CU": "🇨🇺", // Cuba
		"CY": "🇨🇾", // Cyprus
		"CZ": "🇨🇿", // Czech Republic
		"DK": "🇩🇰", // Denmark
		"DO": "🇩🇴", // Dominican Republic
		"EC": "🇪🇨", // Ecuador
		"EG": "🇪🇬", // Egypt
		"SV": "🇸🇻", // El Salvador
		"EE": "🇪🇪", // Estonia
		"FI": "🇫🇮", // Finland
		"FR": "🇫🇷", // France
		"GE": "🇬🇪", // Georgia
		"DE": "🇩🇪", // Germany
		"GH": "🇬🇭", // Ghana
		"GB": "🇬🇧", // Great Britain (United Kingdom)
		"GR": "🇬🇷", // Greece
		"GT": "🇬🇹", // Guatemala
		"HK": "🇭🇰", // Hong Kong
		"HN": "🇭🇳", // Honduras
		"HU": "🇭🇺", // Hungary
		"IS": "🇮🇸", // Iceland
		"IN": "🇮🇳", // India
		"ID": "🇮🇩", // Indonesia
		"IR": "🇮🇷", // Iran
		"IQ": "🇮🇶", // Iraq
		"IE": "🇮🇪", // Ireland
		"IL": "🇮🇱", // Israel
		"IT": "🇮🇹", // Italy
		"JM": "🇯🇲", // Jamaica
		"JP": "🇯🇵", // Japan
		"JO": "🇯🇴", // Jordan
		"KZ": "🇰🇿", // Kazakhstan
		"KE": "🇰🇪", // Kenya
		"KR": "🇰🇷", // Korea, South
		"KW": "🇰🇼", // Kuwait
		"LV": "🇱🇻", // Latvia
		"LB": "🇱🇧", // Lebanon
		"LY": "🇱🇾", // Libya
		"LT": "🇱🇹", // Lithuania
		"LU": "🇱🇺", // Luxembourg
		"MK": "🇲🇰", // North Macedonia
		"MY": "🇲🇾", // Malaysia
		"MT": "🇲🇹", // Malta
		"MX": "🇲🇽", // Mexico
		"MD": "🇲🇩", // Moldova
		"MN": "🇲🇳", // Mongolia
		"ME": "🇲🇪", // Montenegro
		"MA": "🇲🇦", // Morocco
		"MM": "🇲🇲", // Myanmar
		"NP": "🇳🇵", // Nepal
		"NL": "🇳🇱", // Netherlands
		"NZ": "🇳🇿", // New Zealand
		"NI": "🇳🇮", // Nicaragua
		"NG": "🇳🇬", // Nigeria
		"NO": "🇳🇴", // Norway
		"OM": "🇴🇲", // Oman
		"PK": "🇵🇰", // Pakistan
		"PA": "🇵🇦", // Panama
		"PY": "🇵🇾", // Paraguay
		"PE": "🇵🇪", // Peru
		"PH": "🇵🇭", // Philippines
		"PL": "🇵🇱", // Poland
		"PT": "🇵🇹", // Portugal
		"PR": "🇵🇷", // Puerto Rico
		"QA": "🇶🇦", // Qatar
		"RO": "🇷🇴", // Romania
		"RU": "🇷🇺", // Russia
		"SA": "🇸🇦", // Saudi Arabia
		"RS": "🇷🇸", // Serbia
		"SG": "🇸🇬", // Singapore
		"SK": "🇸🇰", // Slovakia
		"SI": "🇸🇮", // Slovenia
		"ZA": "🇿🇦", // South Africa
		"ES": "🇪🇸", // Spain
		"LK": "🇱🇰", // Sri Lanka
		"SD": "🇸🇩", // Sudan
		"SE": "🇸🇪", // Sweden
		"CH": "🇨🇭", // Switzerland
		"SY": "🇸🇾", // Syria
		"TW": "🇹🇼", // Taiwan
		"TJ": "🇹🇯", // Tajikistan
		"TH": "🇹🇭", // Thailand
		"TN": "🇹🇳", // Tunisia
		"TR": "🇹🇷", // Turkey
		"TM": "🇹🇲", // Turkmenistan
		"UA": "🇺🇦", // Ukraine
		"AE": "🇦🇪", // United Arab Emirates
		"US": "🇺🇸", // United States
		"UY": "🇺🇾", // Uruguay
		"UZ": "🇺🇿", // Uzbekistan
		"VE": "🇻🇪", // Venezuela
		"VN": "🇻🇳", // Vietnam
	}

	// Return the flag emoji if it exists in the map, otherwise return the country code
	if flag, ok := countryFlags[countryCode]; ok {
		return flag
	}

	// If no flag is found, create a Unicode flag from the country code
	// This uses the Regional Indicator Symbol Letters which create flags when paired
	// Each letter A-Z is represented by a Unicode code point from U+1F1E6 to U+1F1FF
	if len(countryCode) == 2 {
		// Convert each letter to its regional indicator symbol
		first := 127462 + int(countryCode[0]) - 'A' // 127462 is the Unicode point for 🇦 (Regional Indicator A)
		second := 127462 + int(countryCode[1]) - 'A'
		return string(rune(first)) + string(rune(second))
	}

	return "🌍" // Default earth globe if conversion isn't possible
}

func (vpn *DefaultVpn) Connect(cfg string) error {
	installationToken, err := vpn.cfgProvider.GetInstallationToken()
	if err != nil {
		return errors.NewConfigError("get installation token", err)
	}

	server, err := vpn.serverStorage.GetByServerName(cfg)
	if err != nil {
		return errors.NewConfigError("get server by name", err)
	}

	if server == nil {
		return errors.NewUserError(fmt.Sprintf("Server %s not found", cfg), errors.ErrNotFound)
	}

	ifaceName := sanitizeInterfaceName(server.Hostname)

	output.PrintMsg(fmt.Sprintf("Connecting to %s...", ifaceName), output.MsgOutput)

	keyData, err := vpn.cfgProvider.Get()
	if err != nil || keyData.PrivateKey == "" || keyData.PublicKey == "" {
		// Generate new keys if none exist or there was an error
		publicKey, _, privateKey, keyErr := generateKeys()
		if keyErr != nil {
			return errors.NewVPNError("generate keys", keyErr)
		}

		err = vpn.cfgProvider.StoreData(publicKey.String(), privateKey.String())
		if err != nil {
			return errors.NewConfigError("store keys", err)
		}

		keyData, err = vpn.cfgProvider.Get()
		if err != nil {
			return errors.NewConfigError("get updated config", err)
		}
	}

	ipAddrs, err := vpn.holocron.VpnRegisterPublicKey(installationToken, keyData.PublicKey)
	if err != nil {
		return errors.NewNetworkError("register public key", err)
	}

	// Create WireGuard interface
	if err := vpn.wgManager.CreateInterface(ifaceName); err != nil {
		// If interface already exists, try to remove it first
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("Interface %s already exists, removing and recreating", ifaceName)
			if rmErr := vpn.wgManager.RemoveInterface(ifaceName); rmErr != nil {
				return errors.NewVPNError("remove existing interface", rmErr)
			}
			if err := vpn.wgManager.CreateInterface(ifaceName); err != nil {
				return errors.NewVPNError("create interface", err)
			}
		} else {
			return errors.NewVPNError("create interface", err)
		}
	}

	// Determine port from server's port ranges, fallback to 51820
	port := 51820
	if len(server.PortRanges) > 0 {
		port = server.PortRanges[0].From
	}

	// Configure WireGuard
	wgCfg := wireguard.DeviceConfig{
		PrivateKey: keyData.PrivateKey,
		Peers: []wireguard.PeerConfig{{
			PublicKey:  server.PublicKey,
			Endpoint:   fmt.Sprintf("%s:%d", server.IPv4AddrIn, port),
			AllowedIPs: []string{"0.0.0.0/0", "::/0"},
		}},
	}
	if err := vpn.wgManager.Configure(ifaceName, wgCfg); err != nil {
		vpn.wgManager.RemoveInterface(ifaceName) // Cleanup on failure
		return errors.NewVPNError("configure device", err)
	}

	// Assign addresses
	addrs := []string{ipAddrs.IpV4, ipAddrs.IpV6}
	if err := vpn.wgManager.AssignAddresses(ifaceName, addrs); err != nil {
		vpn.wgManager.RemoveInterface(ifaceName)
		return errors.NewVPNError("assign addresses", err)
	}

	// Bring interface up
	if err := vpn.wgManager.SetInterfaceUp(ifaceName); err != nil {
		vpn.wgManager.RemoveInterface(ifaceName)
		return errors.NewVPNError("set interface up", err)
	}

	// Parse endpoint IP for routing
	endpointIP := net.ParseIP(server.IPv4AddrIn)

	// Add routes
	if err := vpn.wgManager.AddRoutes(ifaceName, []string{"0.0.0.0/0", "::/0"}, endpointIP); err != nil {
		vpn.wgManager.RemoveInterface(ifaceName)
		return errors.NewVPNError("add routes", err)
	}

	output.PrintMsg("Connected.", output.MsgSuccess)
	return nil
}

// sanitizeInterfaceName converts a server hostname to a valid interface name
// Linux interface names are limited to 15 characters
func sanitizeInterfaceName(hostname string) string {
	// Extract just the server name part (before the first dot)
	name := strings.Split(hostname, ".")[0]

	// Replace any non-alphanumeric characters with empty string
	var sanitized strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			sanitized.WriteRune(r)
		}
	}

	result := sanitized.String()

	// Truncate to 15 characters (Linux interface name limit)
	if len(result) > 15 {
		result = result[:15]
	}

	return result
}

func (vpn *DefaultVpn) Disconnect(cfg string) error {
	if cfg == "" {
		devices, err := vpn.wgManager.ListDevices()
		if err != nil {
			return errors.NewVPNError("list devices", err)
		}

		if len(devices) == 0 {
			output.PrintMsg("No active connections.", output.MsgOutput)
			return nil
		}

		for _, name := range devices {
			output.PrintMsg(fmt.Sprintf("Disconnecting from %s...", name), output.MsgOutput)
			if err := vpn.wgManager.RemoveInterface(name); err != nil {
				log.Errorf("Failed to disconnect from %s: %v", name, err)
			} else {
				output.PrintMsg(fmt.Sprintf("Disconnected from %s.", name), output.MsgSuccess)
			}
		}
		return nil
	}

	// If a specific server name is provided, find and disconnect it
	server, err := vpn.serverStorage.GetByServerName(cfg)
	if err != nil {
		return errors.NewConfigError("get server by name", err)
	}

	var ifaceName string
	if server != nil {
		ifaceName = sanitizeInterfaceName(server.Hostname)
	} else {
		// Maybe they provided the interface name directly
		ifaceName = sanitizeInterfaceName(cfg)
	}

	output.PrintMsg(fmt.Sprintf("Disconnecting from %s...", ifaceName), output.MsgOutput)

	if err := vpn.wgManager.RemoveInterface(ifaceName); err != nil {
		return errors.NewVPNError("disconnect", err)
	}

	output.PrintMsg("Disconnected.", output.MsgSuccess)
	return nil
}

func (vpn *DefaultVpn) Status() error {
	devices, err := vpn.wgManager.ListDevices()
	if err != nil {
		return errors.NewVPNError("get status", err)
	}

	if len(devices) == 0 {
		output.PrintMsg("No active connections.", output.MsgOutput)
	} else {
		for _, name := range devices {
			dev, err := vpn.wgManager.GetDevice(name)
			if err != nil {
				log.Errorf("Failed to get device %s: %v", name, err)
				continue
			}
			output.PrintMsg(fmt.Sprintf("Connected to: %s", name), output.MsgSuccess)
			if dev != nil && len(dev.Peers) > 0 {
				for _, peer := range dev.Peers {
					if peer.Endpoint != "" {
						output.PrintMsg(fmt.Sprintf("  Endpoint: %s", peer.Endpoint), output.MsgOutput)
					}
					if !peer.LastHandshakeTime.IsZero() {
						output.PrintMsg(fmt.Sprintf("  Last handshake: %s", peer.LastHandshakeTime.Format("2006-01-02 15:04:05")), output.MsgOutput)
					}
					output.PrintMsg(fmt.Sprintf("  Transfer: ↓ %s / ↑ %s", formatBytes(peer.ReceiveBytes), formatBytes(peer.TransmitBytes)), output.MsgOutput)
				}
			}
		}
	}

	network, err := vpn.holocron.GetVpnNetworkDetails()
	if err != nil {
		return errors.NewNetworkError("get network details", err)
	}

	output.PrintMsg(fmt.Sprintf("	IP Address: %s", network.Ip), output.MsgOutput)
	output.PrintMsg(fmt.Sprintf("	Country: %s", network.Geo.Country), output.MsgOutput)
	output.PrintMsg(fmt.Sprintf("	City: %s", network.Geo.City), output.MsgOutput)

	return nil
}

// formatBytes formats bytes into a human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func generateKeys() (wgtypes.Key, wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, wgtypes.Key{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Get the corresponding public key
	publicKey := privateKey.PublicKey()

	// Generate a pre-shared key if needed
	preSharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, wgtypes.Key{}, fmt.Errorf("failed to generate pre-shared key: %w", err)
	}

	return publicKey, preSharedKey, privateKey, nil
}
