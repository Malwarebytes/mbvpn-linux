package wireguard

import (
	"net"
	"time"
)

// Manager defines the interface for WireGuard operations
type Manager interface {
	// CreateInterface creates a new WireGuard interface
	CreateInterface(name string) error

	// RemoveInterface removes a WireGuard interface
	RemoveInterface(name string) error

	// Configure applies WireGuard configuration to an interface
	Configure(name string, cfg DeviceConfig) error

	// AssignAddresses assigns IP addresses to an interface
	AssignAddresses(name string, addrs []string) error

	// SetInterfaceUp brings the interface up
	SetInterfaceUp(name string) error

	// AddRoutes adds routes for allowed IPs through the interface
	AddRoutes(name string, allowedIPs []string, endpointIP net.IP) error

	// GetDevice returns information about a WireGuard device
	GetDevice(name string) (*Device, error)

	// ListDevices returns a list of all WireGuard interface names
	ListDevices() ([]string, error)

	// IsInterfaceUp checks if an interface is up
	IsInterfaceUp(name string) (bool, error)

	// Close releases any resources held by the manager
	Close() error
}

// DeviceConfig represents WireGuard device configuration
type DeviceConfig struct {
	PrivateKey string
	ListenPort int
	Peers      []PeerConfig
}

// PeerConfig represents a WireGuard peer configuration
type PeerConfig struct {
	PublicKey           string
	Endpoint            string   // "host:port"
	AllowedIPs          []string // CIDR notation: ["0.0.0.0/0", "::/0"]
	PersistentKeepalive int      // seconds, 0 to disable
}

// Device represents a WireGuard device's current state
type Device struct {
	Name       string
	PublicKey  string
	ListenPort int
	Peers      []Peer
}

// Peer represents a WireGuard peer's current state
type Peer struct {
	PublicKey                   string
	Endpoint                    string
	AllowedIPs                  []string
	LastHandshakeTime           time.Time
	TransmitBytes               int64
	ReceiveBytes                int64
	PersistentKeepaliveInterval int
}
