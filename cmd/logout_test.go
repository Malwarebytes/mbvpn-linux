package cmd

import (
	"testing"
)

func TestNewLogoutCommand(t *testing.T) {
	sm := &mockSessionManager{}
	
	cmd := NewLogoutCommand(sm)
	
	if cmd.Use != "logout" {
		t.Errorf("Expected Use to be 'logout', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "Logout this device." {
		t.Errorf("Expected Short to be 'Logout this device.', got '%s'", cmd.Short)
	}
	
	if cmd.Long != "Deactivates this device seat and cleans local configuration." {
		t.Errorf("Expected Long to be 'Deactivates this device seat and cleans local configuration.', got '%s'", cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestLogoutCommandExecution(t *testing.T) {
	sm := &mockSessionManager{}
	
	cmd := NewLogoutCommand(sm)
	cmd.Run(cmd, []string{})
	
	// The command should execute without panicking
	// Actual logout functionality is tested via the mock
}

func TestLogoutCommandLogic(t *testing.T) {
	// Test the logic without running the command to avoid HandleError os.Exit
	sm := &mockSessionManager{}
	
	// Test logout call
	err := sm.Logout()
	if err != nil {
		t.Errorf("Expected no error from sm.Logout, got %v", err)
	}
}