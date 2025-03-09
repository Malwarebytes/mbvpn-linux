package vpn

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Vpn interface {
	Servers()
	Up(cfg string)
	Down(cfg string)
	Status()
}

type DefaultVpn struct {
	cfgProvider config.ConfigProvider
	holocron    remote.Holocron
}

func NewDefaultVpn(cfgProvider config.ConfigProvider, holocron remote.Holocron) Vpn {
	return &DefaultVpn{
		cfgProvider: cfgProvider,
		holocron:    holocron,
	}
}

func (vpn *DefaultVpn) Servers() {
	installationToken, err := vpn.cfgProvider.GetInstallationToken()
	if err != nil {
		log.Panic(err)
	}

	publicKey, _, privateKey, _ := generateKeys()

	ipAddrs, err := vpn.holocron.VpnRegisterPublicKey(installationToken, publicKey.String())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Fetching servers...")

	locations, err := vpn.holocron.GetVpnLocations()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Creating VPN configurations...")

	cmd := exec.Command("sudo", "echo", "Permission granted.")
	err = cmd.Run()
	if err != nil {
		log.Panic(err)
	}

	for _, country := range locations.Countries {
		fmt.Printf("%s:\n", country.Name)
		for _, city := range country.Cities {
			fmt.Printf("  %s:\n", city.Name)
			for i, server := range city.Servers {
				cfgName := fmt.Sprintf("mb-%s-%d", city.Code, i)

				if config.Debug() {
					log.Printf("Creating config: %s\n", cfgName)
				}

				err := writeConfig(cfgName, server, privateKey.String(), ipAddrs.IpV4, ipAddrs.IpV6)
				if err != nil {
					log.Panic(err)
				} else {
					if config.Debug() {
						log.Printf("Config created: %s\n", cfgName)
					}
					fmt.Printf("    %s\n", cfgName)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println(`Call "mbvpn up <server>" to connect and "mbvpn down <server>" to disconnect.`)
}

func (vpn *DefaultVpn) Up(cfg string) {
	fmt.Printf("Connection to %s...\n", cfg)

	cmd := exec.Command("sudo", "wg-quick", "up", cfg)
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

		cmd := exec.Command("sudo", "wg-quick", "down", cfg)
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
		fmt.Println("Connected to servers:")
		for _, s := range servers {
			fmt.Printf("\tConnected to: %s\n", s)
		}
	}
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

func writeConfig(cfgName string, server remote.Server, privateKey string, ipv4 string, ipv6 string) error {
	configPath := filepath.Join("/etc/wireguard", fmt.Sprintf("%s.conf", cfgName))

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

	cmd := exec.Command("sudo", "tee", configPath)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		if _, err := io.WriteString(stdin, content); err != nil {
			log.Fatal("Error writing data to stdin:", err)
		}
	}()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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
