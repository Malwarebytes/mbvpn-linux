package cmd

import (
	"testing"
)

func TestNewLoginCommand(t *testing.T) {
	sm := &mockSessionManager{}

	cmd := NewLoginCommand(sm)

	if cmd.Use != "login" {
		t.Errorf("Expected Use to be 'login', got '%s'", cmd.Use)
	}

	if cmd.Short != "Activates this device with your Malwarebytes license." {
		t.Errorf("Expected Short to be 'Activates this device with your Malwarebytes license.', got '%s'", cmd.Short)
	}

	expectedLong := `Uses provided credentials to perform activation for this device.
   The command claims available seat of your Malwarebytes license.`
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, cmd.Long)
	}

	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestLoginCommandFlags(t *testing.T) {
	sm := &mockSessionManager{}
	cmd := NewLoginCommand(sm)
	mbCodeFlag := cmd.Flags().Lookup("code")
	if mbCodeFlag == nil {
		t.Error("Expected code flag to be defined")
	} else {
		if mbCodeFlag.Shorthand != "c" {
			t.Errorf("Expected code shorthand to be 'c', got '%s'", mbCodeFlag.Shorthand)
		}
		if mbCodeFlag.Usage != "MB-code." {
			t.Errorf("Expected code usage to be 'MB-code.', got '%s'", mbCodeFlag.Usage)
		}
	}
}

func TestLoginCommandWithMBCode(t *testing.T) {
	sm := &mockSessionManager{}

	cmd := NewLoginCommand(sm)
	cmd.Flags().Set("code", "test-mb-code")
	cmd.Run(cmd, []string{})

	// Verify that LoginWithCode was called with the correct parameter
	// Note: We can't directly verify this with our simple mock, but we can test the logic
}

func TestLoginCommandLogic(t *testing.T) {
	sm := &mockSessionManager{}
	err := sm.LoginWithCode("test-mb-code")
	if err != nil {
		t.Errorf("Expected no error from sm.LoginWithCode, got %v", err)
	}
}
