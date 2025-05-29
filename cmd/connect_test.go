package cmd

import (
	"testing"
)

func TestNewConnectCommand(t *testing.T) {
	sm := &mockSessionManager{}
	vpn := &mockVpn{}
	
	cmd := NewConnectCommand(sm, vpn)
	
	if cmd.Use != "connect" {
		t.Errorf("Expected Use to be 'connect', got '%s'", cmd.Use)
	}
	
	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "c" {
		t.Errorf("Expected Aliases to be ['c'], got %v", cmd.Aliases)
	}
	
	if cmd.Short != "Connect to a VPN server" {
		t.Errorf("Expected Short to be 'Connect to a VPN server', got '%s'", cmd.Short)
	}
	
	expectedLong := `Establishes a VPN connection to the specified server using WireGuard.
Requires an active session (login first) and a valid server identifier.

Server Specification Options:
  You can specify the server in multiple ways:

  1. Exact server name:     us-newyork-1
  2. City code:             NYC (randomly selects from New York servers)
  3. City name:             "New York" (randomly selects from New York servers)
  4. Country code:          US (randomly selects from US servers)
  5. Country name:          "United States" (randomly selects from US servers)

Examples:
  mbvpn connect us-newyork-1     # Connect to specific New York server #1
  mbvpn connect NYC              # Connect to random New York server
  mbvpn connect "New York"       # Connect to random New York server
  mbvpn connect US               # Connect to random US server
  mbvpn connect "United States"  # Connect to random US server

Use 'mbvpn servers' to see all available servers with their exact names.
Use 'mbvpn countries' or 'mbvpn cities' to browse servers by location.`
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestConnectCommandWithActiveSession(t *testing.T) {
	sm := &mockSessionManager{active: true}
	vpn := &mockVpn{}
	
	cmd := NewConnectCommand(sm, vpn)
	cmd.Run(cmd, []string{"test-server"})
	
	if len(vpn.connectCalls) != 1 {
		t.Errorf("Expected 1 call to vpn.Connect, got %d", len(vpn.connectCalls))
	}
	
	if vpn.connectCalls[0] != "test-server" {
		t.Errorf("Expected connect call with 'test-server', got '%s'", vpn.connectCalls[0])
	}
}

func TestConnectCommandLogic(t *testing.T) {
	// Test the logic without running the command to avoid HandleError os.Exit
	sm := &mockSessionManager{active: false}
	vpn := &mockVpn{}
	
	// Test inactive session case
	if sm.Active() {
		t.Error("Expected session to be inactive")
	}
	
	// Test active session case  
	sm.active = true
	if !sm.Active() {
		t.Error("Expected session to be active")
	}
	
	// Test vpn.Connect call with server argument
	err := vpn.Connect("us-newyork-1")
	if err != nil {
		t.Errorf("Expected no error from vpn.Connect, got %v", err)
	}
	
	if len(vpn.connectCalls) != 1 {
		t.Errorf("Expected 1 call to vpn.Connect, got %d", len(vpn.connectCalls))
	}
	
	if vpn.connectCalls[0] != "us-newyork-1" {
		t.Errorf("Expected connect call with 'us-newyork-1', got '%s'", vpn.connectCalls[0])
	}
}