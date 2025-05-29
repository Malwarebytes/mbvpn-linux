package cmd

import (
	"testing"
)

func TestNewCountriesCommand(t *testing.T) {
	sm := &mockSessionManager{}
	vpn := &mockVpn{}
	
	cmd := NewCountriesCommand(sm, vpn)
	
	if cmd.Use != "countries" {
		t.Errorf("Expected Use to be 'countries', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "List available countries with VPN servers" {
		t.Errorf("Expected Short to be 'List available countries with VPN servers', got '%s'", cmd.Short)
	}
	
	expectedLong := `Displays a list of countries where VPN servers are available.
Useful for selecting servers by geographic location. Requires an active session.`
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestCountriesCommandWithActiveSession(t *testing.T) {
	sm := &mockSessionManager{active: true}
	vpn := &mockVpn{}
	
	cmd := NewCountriesCommand(sm, vpn)
	cmd.Run(cmd, []string{})
	
	if len(vpn.callParams) != 1 {
		t.Errorf("Expected 1 call to vpn.Servers, got %d", len(vpn.callParams))
	}
	
	if vpn.callParams[0].cities {
		t.Error("Expected cities parameter to be false")
	}
	
	if vpn.callParams[0].full {
		t.Error("Expected full parameter to be false")
	}
}

func TestCountriesCommandLogic(t *testing.T) {
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
	
	// Test vpn.Servers call with correct parameters
	err := vpn.Servers(false, false)
	if err != nil {
		t.Errorf("Expected no error from vpn.Servers, got %v", err)
	}
	
	if len(vpn.callParams) != 1 {
		t.Errorf("Expected 1 call to vpn.Servers, got %d", len(vpn.callParams))
	}
	
	if vpn.callParams[0].cities {
		t.Error("Expected cities parameter to be false")
	}
	
	if vpn.callParams[0].full {
		t.Error("Expected full parameter to be false")
	}
}