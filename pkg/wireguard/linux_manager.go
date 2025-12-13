package wireguard

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// LinuxManager implements Manager using netlink and wgctrl
type LinuxManager struct {
	wgClient *wgctrl.Client
}

// NewManager creates a new WireGuard manager for Linux
func NewManager() (*LinuxManager, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create wgctrl client: %w", err)
	}
	return &LinuxManager{wgClient: client}, nil
}

// Close releases resources held by the manager
func (m *LinuxManager) Close() error {
	if m.wgClient != nil {
		return m.wgClient.Close()
	}
	return nil
}

// CreateInterface creates a new WireGuard interface
func (m *LinuxManager) CreateInterface(name string) error {
	la := netlink.NewLinkAttrs()
	la.Name = name

	wg := &netlink.Wireguard{LinkAttrs: la}
	err := netlink.LinkAdd(wg)
	if err != nil {
		// Check if interface already exists
		if errors.Is(err, syscall.EEXIST) {
			return fmt.Errorf("interface %s already exists", name)
		}
		return fmt.Errorf("failed to create interface %s: %w", name, err)
	}
	return nil
}

// RemoveInterface removes a WireGuard interface
func (m *LinuxManager) RemoveInterface(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		// Interface doesn't exist - not an error
		var linkNotFoundErr netlink.LinkNotFoundError
		if errors.As(err, &linkNotFoundErr) {
			return nil
		}
		return fmt.Errorf("failed to find interface %s: %w", name, err)
	}
	return netlink.LinkDel(link)
}

// Configure applies WireGuard configuration to an interface
func (m *LinuxManager) Configure(name string, cfg DeviceConfig) error {
	// Parse private key
	privKey, err := wgtypes.ParseKey(cfg.PrivateKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	// Build peer configs
	var peers []wgtypes.PeerConfig
	for _, p := range cfg.Peers {
		peerCfg, err := buildPeerConfig(p)
		if err != nil {
			return fmt.Errorf("invalid peer config: %w", err)
		}
		peers = append(peers, peerCfg)
	}

	// Build device config
	wgCfg := wgtypes.Config{
		PrivateKey:   &privKey,
		ReplacePeers: true,
		Peers:        peers,
	}

	// Set listen port if specified
	if cfg.ListenPort > 0 {
		wgCfg.ListenPort = &cfg.ListenPort
	}

	// Apply configuration
	if err := m.wgClient.ConfigureDevice(name, wgCfg); err != nil {
		return fmt.Errorf("failed to configure device %s: %w", name, err)
	}

	return nil
}

func buildPeerConfig(p PeerConfig) (wgtypes.PeerConfig, error) {
	pubKey, err := wgtypes.ParseKey(p.PublicKey)
	if err != nil {
		return wgtypes.PeerConfig{}, fmt.Errorf("invalid public key: %w", err)
	}

	var endpoint *net.UDPAddr
	if p.Endpoint != "" {
		endpoint, err = net.ResolveUDPAddr("udp", p.Endpoint)
		if err != nil {
			return wgtypes.PeerConfig{}, fmt.Errorf("invalid endpoint %s: %w", p.Endpoint, err)
		}
	}

	var allowedIPs []net.IPNet
	for _, cidr := range p.AllowedIPs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return wgtypes.PeerConfig{}, fmt.Errorf("invalid allowed IP %s: %w", cidr, err)
		}
		allowedIPs = append(allowedIPs, *ipnet)
	}

	peerCfg := wgtypes.PeerConfig{
		PublicKey:         pubKey,
		Endpoint:          endpoint,
		AllowedIPs:        allowedIPs,
		ReplaceAllowedIPs: true,
	}

	if p.PersistentKeepalive > 0 {
		keepalive := time.Duration(p.PersistentKeepalive) * time.Second
		peerCfg.PersistentKeepaliveInterval = &keepalive
	}

	return peerCfg, nil
}

// AssignAddresses assigns IP addresses to an interface
func (m *LinuxManager) AssignAddresses(name string, addrs []string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to find interface %s: %w", name, err)
	}

	for _, addr := range addrs {
		// Ensure address is in CIDR notation
		cidrAddr := ensureCIDR(addr)
		nlAddr, err := netlink.ParseAddr(cidrAddr)
		if err != nil {
			return fmt.Errorf("failed to parse address %s: %w", cidrAddr, err)
		}
		if err := netlink.AddrAdd(link, nlAddr); err != nil {
			// Ignore if address already exists
			if !errors.Is(err, syscall.EEXIST) {
				return fmt.Errorf("failed to add address %s: %w", cidrAddr, err)
			}
		}
	}
	return nil
}

// ensureCIDR adds a CIDR suffix if the address doesn't have one
func ensureCIDR(addr string) string {
	// Skip empty addresses
	if addr == "" {
		return addr
	}

	// Already has CIDR notation
	if strings.Contains(addr, "/") {
		return addr
	}

	// Parse the IP to determine if it's IPv4 or IPv6
	ip := net.ParseIP(addr)
	if ip == nil {
		return addr // Return as-is if not a valid IP
	}

	if ip.To4() != nil {
		return addr + "/32"
	}
	return addr + "/128"
}

// SetInterfaceUp brings the interface up
func (m *LinuxManager) SetInterfaceUp(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to find interface %s: %w", name, err)
	}
	return netlink.LinkSetUp(link)
}

// AddRoutes adds routes for allowed IPs through the interface
// endpointIP is used to add a specific route to the VPN server via the original gateway
func (m *LinuxManager) AddRoutes(name string, allowedIPs []string, endpointIP net.IP) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to find interface %s: %w", name, err)
	}

	// Get current default gateway for routing VPN endpoint traffic
	var defaultGW net.IP
	var defaultLinkIndex int
	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return fmt.Errorf("failed to list routes: %w", err)
	}
	for _, r := range routes {
		// Default route can be Dst == nil OR Dst == 0.0.0.0/0
		isDefault := r.Dst == nil
		if !isDefault && r.Dst != nil {
			ones, bits := r.Dst.Mask.Size()
			isDefault = ones == 0 && bits == 32
		}
		if isDefault && r.Gw != nil {
			defaultGW = r.Gw
			defaultLinkIndex = r.LinkIndex
			break
		}
	}

	// Add route to VPN endpoint via original gateway (prevents routing loop)
	if endpointIP != nil && defaultGW != nil {
		var mask net.IPMask
		if endpointIP.To4() != nil {
			mask = net.CIDRMask(32, 32)
		} else {
			mask = net.CIDRMask(128, 128)
		}
		endpointRoute := &netlink.Route{
			LinkIndex: defaultLinkIndex,
			Dst:       &net.IPNet{IP: endpointIP, Mask: mask},
			Gw:        defaultGW,
		}
		if err := netlink.RouteAdd(endpointRoute); err != nil {
			// Ignore if route already exists
			if !errors.Is(err, syscall.EEXIST) {
				return fmt.Errorf("failed to add endpoint route: %w", err)
			}
		}
	}

	// Add routes for allowed IPs through the WireGuard interface
	for _, cidr := range allowedIPs {
		_, dst, err := net.ParseCIDR(cidr)
		if err != nil {
			return fmt.Errorf("failed to parse CIDR %s: %w", cidr, err)
		}

		// For default routes (0.0.0.0/0 or ::/0), we need special handling
		// Using two /1 routes avoids replacing the system default route
		if isDefaultRoute(dst) {
			if err := addSplitDefaultRoutes(link.Attrs().Index, dst); err != nil {
				return err
			}
			continue
		}

		route := &netlink.Route{
			LinkIndex: link.Attrs().Index,
			Dst:       dst,
		}
		if err := netlink.RouteAdd(route); err != nil {
			if !errors.Is(err, syscall.EEXIST) {
				return fmt.Errorf("failed to add route for %s: %w", cidr, err)
			}
		}
	}
	return nil
}

// isDefaultRoute checks if a route is a default route (0.0.0.0/0 or ::/0)
func isDefaultRoute(ipnet *net.IPNet) bool {
	ones, bits := ipnet.Mask.Size()
	return ones == 0 && (bits == 32 || bits == 128)
}

// addSplitDefaultRoutes adds two /1 routes that effectively cover all traffic
// without replacing the system's default route
func addSplitDefaultRoutes(linkIndex int, dst *net.IPNet) error {
	_, bits := dst.Mask.Size()

	if bits == 32 {
		// IPv4: Add 0.0.0.0/1 and 128.0.0.0/1
		routes := []*netlink.Route{
			{
				LinkIndex: linkIndex,
				Dst:       &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(1, 32)},
			},
			{
				LinkIndex: linkIndex,
				Dst:       &net.IPNet{IP: net.IPv4(128, 0, 0, 0), Mask: net.CIDRMask(1, 32)},
			},
		}
		for _, route := range routes {
			if err := netlink.RouteAdd(route); err != nil {
				if !errors.Is(err, syscall.EEXIST) {
					return fmt.Errorf("failed to add split route: %w", err)
				}
			}
		}
	} else {
		// IPv6: Add ::/1 and 8000::/1
		routes := []*netlink.Route{
			{
				LinkIndex: linkIndex,
				Dst:       &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(1, 128)},
			},
			{
				LinkIndex: linkIndex,
				Dst:       &net.IPNet{IP: net.ParseIP("8000::"), Mask: net.CIDRMask(1, 128)},
			},
		}
		for _, route := range routes {
			if err := netlink.RouteAdd(route); err != nil {
				if !errors.Is(err, syscall.EEXIST) {
					return fmt.Errorf("failed to add split IPv6 route: %w", err)
				}
			}
		}
	}
	return nil
}

// GetDevice returns information about a WireGuard device
func (m *LinuxManager) GetDevice(name string) (*Device, error) {
	dev, err := m.wgClient.Device(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device %s: %w", name, err)
	}

	device := &Device{
		Name:       dev.Name,
		PublicKey:  dev.PublicKey.String(),
		ListenPort: dev.ListenPort,
	}

	for _, p := range dev.Peers {
		var allowedIPs []string
		for _, ip := range p.AllowedIPs {
			allowedIPs = append(allowedIPs, ip.String())
		}

		peer := Peer{
			PublicKey:                   p.PublicKey.String(),
			AllowedIPs:                  allowedIPs,
			LastHandshakeTime:           p.LastHandshakeTime,
			TransmitBytes:               p.TransmitBytes,
			ReceiveBytes:                p.ReceiveBytes,
			PersistentKeepaliveInterval: int(p.PersistentKeepaliveInterval.Seconds()),
		}
		if p.Endpoint != nil {
			peer.Endpoint = p.Endpoint.String()
		}
		device.Peers = append(device.Peers, peer)
	}

	return device, nil
}

// ListDevices returns a list of all WireGuard interface names
func (m *LinuxManager) ListDevices() ([]string, error) {
	devices, err := m.wgClient.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	names := make([]string, len(devices))
	for i, d := range devices {
		names[i] = d.Name
	}
	return names, nil
}

// IsInterfaceUp checks if an interface is up
func (m *LinuxManager) IsInterfaceUp(name string) (bool, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		var linkNotFoundErr netlink.LinkNotFoundError
		if errors.As(err, &linkNotFoundErr) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find interface %s: %w", name, err)
	}
	return link.Attrs().Flags&net.FlagUp != 0, nil
}
