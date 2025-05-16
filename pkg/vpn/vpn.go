package vpn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/console"
	"github.com/Malwarebytes/mbvpn/pkg/errors"
	"github.com/Malwarebytes/mbvpn/pkg/output"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
	"github.com/Malwarebytes/mbvpn/pkg/servers"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Vpn interface {
	Servers(showCities bool, showServers bool) error
	Up(cfg string) error
	Down(cfg string) error
	Status() error
}

type DefaultVpn struct {
	cfgProvider   config.ConfigProvider
	holocron      remote.Holocron
	serverStorage servers.ServerStorage
}

func NewDefaultVpn(cfgProvider config.ConfigProvider, holocron remote.Holocron, serverStorage servers.ServerStorage) Vpn {
	return &DefaultVpn{
		cfgProvider:   cfgProvider,
		holocron:      holocron,
		serverStorage: serverStorage,
	}
}

func (vpn *DefaultVpn) Servers(showCities bool, showServers bool) error {
	locations, err := vpn.holocron.GetVpnLocations()
	if err != nil {
		output.PrintMsg("Failed to get server list.", output.MsgError)
		return errors.NewNetworkError("get server list", err)
	}

	serverStorage := servers.DefaultServerStorage{}
	err = serverStorage.Save(locations)
	if err != nil {
		output.PrintMsg("Failed to save server list.", output.MsgError)
		return errors.NewConfigError("save server list", err)
	}

	for _, country := range locations.Countries {
		output.PrintMsg(fmt.Sprintf("%s, %s", country.Name, country.Code), output.MsgOutput)
		for _, city := range country.Cities {
			if showCities {
				output.PrintMsg(fmt.Sprintf("  %s, %s", city.Name, city.Code), output.MsgOutput)
			}
			for _, s := range city.Servers {
				if showServers {
					output.PrintMsg(fmt.Sprintf("    %s", strings.Split(s.Hostname, ".")[0]), output.MsgOutput)
				}
			}
		}
	}
	return nil
}

func (vpn *DefaultVpn) Up(cfg string) error {
	installationToken, err := vpn.cfgProvider.GetInstallationToken()
	if err != nil {
		return errors.NewConfigError("get installation token", err)
	}

	server, err := vpn.serverStorage.GetByServerName(cfg)
	if err != nil {
		return errors.NewConfigError("get server by name", err)
	}
	
	if server == nil {
		output.PrintMsg(fmt.Sprintf("Server %s not found.", cfg), output.MsgError)
		return errors.NewUserError(fmt.Sprintf("Server %s not found", cfg), errors.ErrNotFound)
	}

	cfgName := strings.Split(server.Hostname, ".")[0]

	output.PrintMsg(fmt.Sprintf("Connecting to %s...", cfgName), output.MsgOutput)

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

	cfgPath, err := writeConfig(cfgName, *server, keyData.PrivateKey, ipAddrs.IpV4, ipAddrs.IpV6)
	if err != nil {
		return errors.NewVPNError("write config", err)
	}

	output.PrintMsg(fmt.Sprintf("Calling 'wg-quick up %s'", cfgPath), output.MsgOutput)
	_, err = console.RunCmd(true, "wg-quick", "up", cfgPath)
	if err != nil {
		output.PrintMsg("Failed to connect.", output.MsgError)
		return errors.NewVPNError("connect", err)
	}
	
	output.PrintMsg("Connected.", output.MsgSuccess)
	return nil
}

func (vpn *DefaultVpn) Down(cfg string) error {
	if cfg == "" {
		servers, err := getConnectedServers()
		if err != nil {
			return errors.NewVPNError("get connected servers", err)
		}
		
		for _, s := range servers {
			err := vpn.Down(s)
			if err != nil {
				// Continue trying to disconnect other servers even if one fails
				log.Errorf("Failed to disconnect from %s: %v", s, err)
			}
		}
		return nil
	} 
	
	output.PrintMsg(fmt.Sprintf("Disconnecting from %s...", cfg), output.MsgOutput)

	cfgDir, err := ensureConfigDir()
	if err != nil {
		return errors.NewConfigError("ensure config directory", err)
	}
	
	_, err = console.RunCmd(true, "wg-quick", "down", filepath.Join(cfgDir, cfg+".conf"))
	if err != nil {
		output.PrintMsg("Failed to disconnect.", output.MsgError)
		return errors.NewVPNError("disconnect", err)
	}
	
	output.PrintMsg("Disconnected.", output.MsgSuccess)
	return nil
}

func (vpn *DefaultVpn) Status() error {
	servers, err := getConnectedServers()
	if err != nil {
		return errors.NewVPNError("get status", err)
	}

	if len(servers) == 0 {
		output.PrintMsg("No active connections.", output.MsgOutput)
	} else {
		for _, s := range servers {
			output.PrintMsg(fmt.Sprintf("Connected to: %s", s), output.MsgSuccess)
		}
	}

	network, err := vpn.holocron.GetVpnNetworkDetails()
	if err != nil {
		output.PrintMsg("Failed to get network details.", output.MsgError)
		return errors.NewNetworkError("get network details", err)
	}

	output.PrintMsg(fmt.Sprintf("IP Address: %s", network.Ip), output.MsgOutput)
	output.PrintMsg(fmt.Sprintf("VPN enabled: %t", network.VpnEnabled), output.MsgOutput)
	output.PrintMsg(fmt.Sprintf("Country: %s", network.Geo.Country), output.MsgOutput)
	output.PrintMsg(fmt.Sprintf("City: %s", network.Geo.City), output.MsgOutput)
	
	return nil
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

func writeConfig(cfgName string, server remote.Server, privateKey string, ipv4 string, ipv6 string) (string, error) {
	content := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s, %s

[Peer]
PublicKey = %s
Endpoint = %s:51820
AllowedIPs = 0.0.0.0/0, ::/0`,
		privateKey,
		ipv4,
		ipv6,
		server.PublicKey,
		server.IPv4AddrIn,
	)

	fullPath, err := saveWgConfig(cfgName, content)
	if err != nil {
		return "", fmt.Errorf("failed to save config file: %w", err)
	}

	return fullPath, nil
}

func getConnectedServers() ([]string, error) {	
	output, err := console.RunCmd(true, "wg", "show")
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'wg show': %w", err)
	}

	lines := strings.Split(string(output), "\n")
	n := strings.Count(string(output), "interface:")
	s := make([]string, n)
	i := 0
	for _, l := range lines {
		if strings.HasPrefix(l, "interface") {
			s[i] = strings.TrimPrefix(l, "interface: ")
			i++
		}
	}
	return s, nil
}

func ensureConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "mbvpn", "servers")

	// Create directory with appropriate permissions if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func saveWgConfig(serverName string, configContent string) (string, error) {
	configDir, err := ensureConfigDir()
	if err != nil {
		return "", err
	}

	// Sanitize the server name to avoid path traversal attacks
	serverName = filepath.Base(serverName)

	// Create the file with restricted permissions (600) as it contains private keys
	filePath := filepath.Join(configDir, serverName+".conf")
	return filePath, os.WriteFile(filePath, []byte(configContent), 0o600)
}