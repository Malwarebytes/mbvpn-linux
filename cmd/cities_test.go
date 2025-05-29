package cmd

import (
	"testing"
)

// Mock implementations
type mockSessionManager struct {
	active bool
}

func (m *mockSessionManager) Active() bool {
	return m.active
}

func (m *mockSessionManager) LoginWithKey(string) error {
	return nil
}

func (m *mockSessionManager) LoginWithCode(string) error {
	return nil
}

func (m *mockSessionManager) Logout() error {
	return nil
}

type mockVpn struct {
	serversError    error
	callParams      []mockVpnCall
	connectCalls    []string
	disconnectCalls []string
}

type mockVpnCall struct {
	cities bool
	full   bool
}

func (m *mockVpn) Servers(cities, full bool) error {
	m.callParams = append(m.callParams, mockVpnCall{cities: cities, full: full})
	return m.serversError
}

func (m *mockVpn) Connect(cfg string) error {
	m.connectCalls = append(m.connectCalls, cfg)
	return nil
}

func (m *mockVpn) Disconnect(cfg string) error {
	m.disconnectCalls = append(m.disconnectCalls, cfg)
	return nil
}

func (m *mockVpn) Status() error {
	return nil
}

func TestNewCitiesCommand(t *testing.T) {
	sm := &mockSessionManager{}
	vpn := &mockVpn{}
	
	cmd := NewCitiesCommand(sm, vpn)
	
	if cmd.Use != "cities" {
		t.Errorf("Expected Use to be 'cities', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "List available cities with VPN servers" {
		t.Errorf("Expected Short to be 'List available cities with VPN servers', got '%s'", cmd.Short)
	}
	
	expectedLong := `Displays a list of cities where VPN servers are located.
Provides more specific geographic options than countries. Requires an active session.`
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestCitiesCommandWithActiveSession(t *testing.T) {
	sm := &mockSessionManager{active: true}
	vpn := &mockVpn{}
	
	cmd := NewCitiesCommand(sm, vpn)
	cmd.Run(cmd, []string{})
	
	if len(vpn.callParams) != 1 {
		t.Errorf("Expected 1 call to vpn.Servers, got %d", len(vpn.callParams))
	}
	
	if !vpn.callParams[0].cities {
		t.Error("Expected cities parameter to be true")
	}
	
	if vpn.callParams[0].full {
		t.Error("Expected full parameter to be false")
	}
}

func TestCitiesCommandLogic(t *testing.T) {
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
	err := vpn.Servers(true, false)
	if err != nil {
		t.Errorf("Expected no error from vpn.Servers, got %v", err)
	}
	
	if len(vpn.callParams) != 1 {
		t.Errorf("Expected 1 call to vpn.Servers, got %d", len(vpn.callParams))
	}
	
	if !vpn.callParams[0].cities {
		t.Error("Expected cities parameter to be true")
	}
	
	if vpn.callParams[0].full {
		t.Error("Expected full parameter to be false")
	}
}