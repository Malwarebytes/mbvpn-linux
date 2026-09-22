package daemon

import (
	"crypto/rand"
	"fmt"
	"net"

	"github.com/malwarebytes/mbvpn-linux/pkg/wireguard"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func generateInterfaceName() (string, error) {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("mbv%x", bytes), nil
}

func generateKeys() (wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, err
	}
	return privateKey.PublicKey(), privateKey, nil
}

func connectInterface(manager wireguard.Manager, name, privateKey, peerKey, endpoint string, addresses []string) error {
	if err := manager.CreateInterface(name); err != nil {
		return err
	}
	cleanup := true
	defer func() {
		if cleanup {
			_ = manager.RemoveInterface(name)
		}
	}()
	if err := manager.Configure(name, wireguard.DeviceConfig{PrivateKey: privateKey, Peers: []wireguard.PeerConfig{{
		PublicKey: peerKey, Endpoint: endpoint, AllowedIPs: []string{"0.0.0.0/0", "::/0"},
	}}}); err != nil {
		return err
	}
	if err := manager.AssignAddresses(name, addresses); err != nil {
		return err
	}
	if err := manager.SetInterfaceUp(name); err != nil {
		return err
	}
	host, _, err := net.SplitHostPort(endpoint)
	if err != nil {
		return err
	}
	endpointIP := net.ParseIP(host)
	if endpointIP == nil {
		return fmt.Errorf("invalid endpoint IP")
	}
	if err := manager.AddRoutes(name, []string{"0.0.0.0/0", "::/0"}, endpointIP); err != nil {
		return err
	}
	cleanup = false
	return nil
}
