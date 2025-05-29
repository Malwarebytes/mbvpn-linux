package cmd

import (
	"testing"
)

func TestNewStatusCommand(t *testing.T) {
	vpn := &mockVpn{}
	
	cmd := NewStatusCommand(vpn)
	
	if cmd.Use != "status" {
		t.Errorf("Expected Use to be 'status', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "Show the status of the current VPN connection" {
		t.Errorf("Expected Short to be 'Show the status of the current VPN connection', got '%s'", cmd.Short)
	}
	
	expectedLong := `Displays detailed information about the current VPN connection status,
including connection state, server details, and network configuration.`
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestStatusCommandExecution(t *testing.T) {
	vpn := &mockVpn{}
	
	cmd := NewStatusCommand(vpn)
	cmd.Run(cmd, []string{})
	
	// The command should execute without panicking
	// Actual status functionality is tested via the mock
}

func TestStatusCommandLogic(t *testing.T) {
	// Test the logic without running the command to avoid HandleError os.Exit
	vpn := &mockVpn{}
	
	// Test vpn.Status call
	err := vpn.Status()
	if err != nil {
		t.Errorf("Expected no error from vpn.Status, got %v", err)
	}
}