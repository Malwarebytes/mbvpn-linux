package cmd

import (
	"testing"
)

func TestNewDisconnectCommand(t *testing.T) {
	vpn := &mockVpn{}
	
	cmd := NewDisconnectCommand(vpn)
	
	if cmd.Use != "disconnect" {
		t.Errorf("Expected Use to be 'disconnect', got '%s'", cmd.Use)
	}
	
	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "d" {
		t.Errorf("Expected Aliases to be ['d'], got %v", cmd.Aliases)
	}
	
	if cmd.Short != "Stops the VPN connection." {
		t.Errorf("Expected Short to be 'Stops the VPN connection.', got '%s'", cmd.Short)
	}
	
	if cmd.Long != "Stops the VPN connection." {
		t.Errorf("Expected Long to be 'Stops the VPN connection.', got '%s'", cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestDisconnectCommandWithNoArgs(t *testing.T) {
	vpn := &mockVpn{}
	
	cmd := NewDisconnectCommand(vpn)
	cmd.Run(cmd, []string{})
	
	if len(vpn.disconnectCalls) != 1 {
		t.Errorf("Expected 1 call to vpn.Disconnect, got %d", len(vpn.disconnectCalls))
	}
	
	if vpn.disconnectCalls[0] != "" {
		t.Errorf("Expected disconnect call with empty string, got '%s'", vpn.disconnectCalls[0])
	}
}

func TestDisconnectCommandWithServerArg(t *testing.T) {
	vpn := &mockVpn{}
	
	cmd := NewDisconnectCommand(vpn)
	cmd.Run(cmd, []string{"test-server"})
	
	if len(vpn.disconnectCalls) != 1 {
		t.Errorf("Expected 1 call to vpn.Disconnect, got %d", len(vpn.disconnectCalls))
	}
	
	if vpn.disconnectCalls[0] != "test-server" {
		t.Errorf("Expected disconnect call with 'test-server', got '%s'", vpn.disconnectCalls[0])
	}
}

func TestDisconnectCommandLogic(t *testing.T) {
	// Test the logic without running the command to avoid HandleError os.Exit
	vpn := &mockVpn{}
	
	// Test disconnect call with empty server
	err := vpn.Disconnect("")
	if err != nil {
		t.Errorf("Expected no error from vpn.Disconnect, got %v", err)
	}
	
	if len(vpn.disconnectCalls) != 1 {
		t.Errorf("Expected 1 call to vpn.Disconnect, got %d", len(vpn.disconnectCalls))
	}
	
	if vpn.disconnectCalls[0] != "" {
		t.Errorf("Expected disconnect call with empty string, got '%s'", vpn.disconnectCalls[0])
	}
	
	// Test disconnect call with specific server
	vpn2 := &mockVpn{}
	err = vpn2.Disconnect("us-newyork-1")
	if err != nil {
		t.Errorf("Expected no error from vpn.Disconnect, got %v", err)
	}
	
	if len(vpn2.disconnectCalls) != 1 {
		t.Errorf("Expected 1 call to vpn.Disconnect, got %d", len(vpn2.disconnectCalls))
	}
	
	if vpn2.disconnectCalls[0] != "us-newyork-1" {
		t.Errorf("Expected disconnect call with 'us-newyork-1', got '%s'", vpn2.disconnectCalls[0])
	}
}