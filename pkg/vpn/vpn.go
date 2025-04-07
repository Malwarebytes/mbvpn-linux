package vpn

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
	"github.com/Malwarebytes/mbvpn/pkg/servers"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Vpn interface {
	Servers(showCities bool, showServers bool)
	Up(cfg string)
	Down(cfg string)
	Status()
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

func (vpn *DefaultVpn) Servers(showCities bool, showServers bool) {
	installationToken, err := vpn.cfgProvider.GetInstallationToken()
	if err != nil {
		log.Panic(err)
	}

	m, err := vpn.holocron.CheckDevice(installationToken)
	if err != nil {
		log.Panic(err)
	}
	if m.Status != remote.DeviceStatusLicensed {
		fmt.Println("You need to activate your device first. Use `mbvpn login` to activate.")
		return
	}

	locations, err := vpn.holocron.GetVpnLocations()
	if err != nil {
		log.Panic(err)
	}

	serverStorage := servers.DefaultServerStorage{}
	err = serverStorage.Save(locations)
	if err != nil {
		fmt.Println("Failed to save server list:", err)
		return
	}

	for _, country := range locations.Countries {
		fmt.Printf("%s, %s\n", country.Name, country.Code)
		for _, city := range country.Cities {
			if showCities {
				fmt.Printf("  %s, %s\n", city.Name, city.Code)
			}
			for _, s := range city.Servers {
				if showServers {
					fmt.Printf("    %s\n", strings.Split(s.Hostname, ".")[0])
				}
			}
		}
	}

	fmt.Println()
	fmt.Println(`Call "mbvpn up <server>" to connect and "mbvpn down <server>" to disconnect.`)
}

func (vpn *DefaultVpn) Up(cfg string) {
	installationToken, err := vpn.cfgProvider.GetInstallationToken()
	if err != nil {
		log.Panic(err)
	}


	server, err := vpn.serverStorage.GetByServerName(cfg)
	if err != nil {
		log.Panic(err)
	}
	if server == nil {
		fmt.Printf("Server %s not found.\n", cfg)
		return
	}

	cfgName := strings.Split(server.Hostname, ".")[0]

	fmt.Printf("Connecting to %s...\n", cfgName)

	publicKey, _, privateKey, _ := generateKeys()

	ipAddrs, err := vpn.holocron.VpnRegisterPublicKey(installationToken, publicKey.String())
	if err != nil {
		log.Panic(err)
	}

	cfgPath, err := writeConfig(cfgName, *server, privateKey.String(), ipAddrs.IpV4, ipAddrs.IpV6)
	if err != nil {
		log.Panic(err)
		return
	}

	fmt.Printf("Calling `sudo wg-quick up %s`\n", cfgPath)
	cmd := exec.Command("sudo", "wg-quick", "up", cfgPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Panic(err)
	} else {
		fmt.Println("Connected.")
	}
}

func (vpn *DefaultVpn) Down(cfg string) {
	if cfg == "" {
		servers, err := getConnectedServers()
		if err != nil {
			log.Panic(err)
		} else {
			for _, s := range servers {
				vpn.Down(s)
			}
		}
	} else {
		fmt.Printf("Disconnecting from %s...\n", cfg)

		cfgDir, err := ensureConfigDir()
		if err != nil {
			log.Panic(err)
			return
		}
		cmd := exec.Command("sudo", "wg-quick", "down", filepath.Join(cfgDir, cfg+".conf"))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Disconnected.")
		}
	}
}

func (vpn *DefaultVpn) Status() {
	servers, err := getConnectedServers()
	if err != nil {
		log.Panic(err)
		return
	}

	if len(servers) == 0 {
		fmt.Println("No active connections.")
	} else {
		for _, s := range servers {
			fmt.Printf("Connected to: %s\n", s)
		}
	}

	network, err := vpn.holocron.GetVpnNetworkDetails()
	if err != nil {
		fmt.Println("Failed to get network details.")
		if config.Debug() {
			log.Panic(err)
		}
		return
	}

	fmt.Printf("IP Address: %s\n", network.Ip)
	fmt.Printf("VPN enabled: %t\n", network.VpnEnabled)
	fmt.Printf("Country: %s\n", network.Geo.Country)
	fmt.Printf("City: %s\n", network.Geo.City)
}

func generateKeys() (wgtypes.Key, wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	// Get the corresponding public key
	publicKey := privateKey.PublicKey()

	// Generate a pre-shared key if needed
	preSharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		panic(err)
	}

	return publicKey, preSharedKey, privateKey, nil
}

func intPtr(i int) *int {
	return &i
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
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
		fmt.Println("Failed to save config file:", err)
		return "", err
	}

	return fullPath, nil
}

func getConnectedServers() ([]string, error) {
	cmd := exec.Command("sudo", "wg", "show")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
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
