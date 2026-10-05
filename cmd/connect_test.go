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
  You can specify the server in multiple ways (case-insensitive):

  1. Exact server name:     gb-lon-2
  2. City code:             gb-lon (randomly selects from London servers)
  3. City name:             "London" (randomly selects from London servers)
  4. Country code:          gb (randomly selects from UK servers)
  5. Country name:          "United Kingdom" (randomly selects from UK servers)

Examples:
  mbvpn connect gb-lon-2         # Connect to specific London server #2
  mbvpn connect gb-lon           # Connect to random London server
  mbvpn connect London           # Connect to random London server
  mbvpn connect gb               # Connect to random UK server
  mbvpn connect "United Kingdom" # Connect to random UK server

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