package wireguard

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIfNotPrivileged(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root or CAP_NET_ADMIN")
	}
}

func TestNewManager(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	assert.NotNil(t, mgr.wgClient)
}

func TestCreateAndRemoveInterface(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest0"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// Verify interface exists
	up, err := mgr.IsInterfaceUp(testInterface)
	require.NoError(t, err)
	assert.False(t, up) // Interface is created but not up

	// Remove interface
	err = mgr.RemoveInterface(testInterface)
	require.NoError(t, err)

	// Verify interface is gone
	up, err = mgr.IsInterfaceUp(testInterface)
	require.NoError(t, err)
	assert.False(t, up)
}

func TestCreateInterfaceAlreadyExists(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest1"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)
	defer mgr.RemoveInterface(testInterface)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// Try to create again - should fail
	err = mgr.CreateInterface(testInterface)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRemoveNonExistentInterface(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	// Removing non-existent interface should not error
	err = mgr.RemoveInterface("nonexistent123")
	assert.NoError(t, err)
}

func TestListDevices(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest2"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)
	defer mgr.RemoveInterface(testInterface)

	// Get initial device count
	initialDevices, err := mgr.ListDevices()
	require.NoError(t, err)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// List devices should include our new interface
	devices, err := mgr.ListDevices()
	require.NoError(t, err)
	assert.Len(t, devices, len(initialDevices)+1)
	assert.Contains(t, devices, testInterface)
}

func TestConfigureDevice(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest3"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)
	defer mgr.RemoveInterface(testInterface)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// Configure with a test configuration
	cfg := DeviceConfig{
		PrivateKey: "WG8ZcfSVD/oaEWcUW5lrJsLVWgYA8X0m9Iuv9U3XSWA=", // Test key only
		Peers: []PeerConfig{
			{
				PublicKey:  "bPfJDdgBzYmjXLq0S+VQkzf5GdOOKZ5zjP8HCBXDYV4=", // Test key only
				Endpoint:   "192.0.2.1:51820",
				AllowedIPs: []string{"10.0.0.0/24"},
			},
		},
	}

	err = mgr.Configure(testInterface, cfg)
	require.NoError(t, err)

	// Verify configuration
	dev, err := mgr.GetDevice(testInterface)
	require.NoError(t, err)
	require.NotNil(t, dev)
	assert.Len(t, dev.Peers, 1)
}

func TestAssignAddresses(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest4"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)
	defer mgr.RemoveInterface(testInterface)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// Assign addresses
	err = mgr.AssignAddresses(testInterface, []string{"10.100.100.1/32", "fd00::1/128"})
	require.NoError(t, err)
}

func TestSetInterfaceUp(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest5"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)
	defer mgr.RemoveInterface(testInterface)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// Assign an address first (required for bringing up)
	err = mgr.AssignAddresses(testInterface, []string{"10.100.100.2/32"})
	require.NoError(t, err)

	// Bring interface up
	err = mgr.SetInterfaceUp(testInterface)
	require.NoError(t, err)

	// Verify interface is up
	up, err := mgr.IsInterfaceUp(testInterface)
	require.NoError(t, err)
	assert.True(t, up)
}

func TestInvalidPrivateKey(t *testing.T) {
	skipIfNotPrivileged(t)

	mgr, err := NewManager()
	require.NoError(t, err)
	defer mgr.Close()

	const testInterface = "wgtest6"

	// Ensure clean state
	_ = mgr.RemoveInterface(testInterface)
	defer mgr.RemoveInterface(testInterface)

	// Create interface
	err = mgr.CreateInterface(testInterface)
	require.NoError(t, err)

	// Try to configure with invalid key
	cfg := DeviceConfig{
		PrivateKey: "invalid-key",
	}

	err = mgr.Configure(testInterface, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid private key")
}

func TestIsDefaultRoute(t *testing.T) {
	tests := []struct {
		cidr     string
		expected bool
	}{
		{"0.0.0.0/0", true},
		{"::/0", true},
		{"10.0.0.0/8", false},
		{"192.168.0.0/16", false},
		{"0.0.0.0/1", false},
		{"128.0.0.0/1", false},
	}

	for _, tt := range tests {
		t.Run(tt.cidr, func(t *testing.T) {
			_, ipnet, err := net.ParseCIDR(tt.cidr)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, isDefaultRoute(ipnet))
		})
	}
}
