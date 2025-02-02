package vpn

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/remote"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func Servers() {
	installationToken, _ := config.GetInstallationToken()

	publicKey, _, privateKey, _ := generateKeys()

	ipAddrs, err := remote.VpnRegisterPublicKey(installationToken, publicKey.String())
	if err != nil {
		log.Panic(err)
	}

	locations, err := remote.GetVpnLocations(installationToken)
	if err != nil {
		log.Panic(err)
	}

	for _, country := range locations.Countries {
		fmt.Printf("%s:\n", country.Name)
		for _, city := range country.Cities {
			fmt.Printf("  %s:\n", city.Name)
			for i, server := range city.Servers {
				cfgName := fmt.Sprintf("mb-%s-%d", city.Code, i)

				if config.Debug {
					log.Printf("Creating config: %s\n", cfgName)
				}

				err := writeConfig(cfgName, server, privateKey.String(), ipAddrs.IpV4, ipAddrs.IpV6)
				if err != nil {
					log.Panic(err)
				} else {
					if config.Debug {
						log.Printf("Config created: %s\n", cfgName)
					}
					fmt.Printf("    %s\n", cfgName)
				}
			}
		}
	}

  fmt.Println(`Call "mbvpn up <server>" to connect and "mbvpn down <server>" to disconnect.`)

	// networkDetails, err := remote.GetVpnNetworkDetails(installationToken)
	// if err != nil {
	//   log.Panic(err)
	// }
	//
	// fmt.Printf("Current country: %s", networkDetails.Geo.Country)
	//
	// //Connect
	//
	// fmt.Printf("Current country: %s", networkDetails.Geo.Country)
	//
	// //Disconnect
	//
	// fmt.Printf("Current country: %s", networkDetails.Geo.Country)
}

func Up(cfg string) {
	_, err := os.Stat(cfg)
	if errors.Is(err, os.ErrNotExist) {
    fmt.Fprintf(os.Stderr, "Config %s doesn't exist. Use \"servers\" command to see available configurations.\n", cfg)
    os.Exit(2)
  }

  fmt.Printf("Connection to %s...\n", cfg)

  //TODO
  // _ := config.Debug


  fmt.Println("Connected!")
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
AllowedIPs = 0.0.0.0/0, ::/0
`,
		privateKey,
		ipv4,
		ipv6,
		server.PublicKey,
		server.IPv4AddrIn,
	)

	return os.WriteFile(configPath, []byte(content), 0o600)
}
