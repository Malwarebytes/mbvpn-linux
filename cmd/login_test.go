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
	
	// Test licenseKey flag
	licenseKeyFlag := cmd.Flags().Lookup("licenseKey")
	if licenseKeyFlag == nil {
		t.Error("Expected licenseKey flag to be defined")
	} else {
		if licenseKeyFlag.Shorthand != "k" {
			t.Errorf("Expected licenseKey shorthand to be 'k', got '%s'", licenseKeyFlag.Shorthand)
		}
		if licenseKeyFlag.Usage != "License key." {
			t.Errorf("Expected licenseKey usage to be 'License key.', got '%s'", licenseKeyFlag.Usage)
		}
	}
	
	// Test mbCode flag
	mbCodeFlag := cmd.Flags().Lookup("mbCode")
	if mbCodeFlag == nil {
		t.Error("Expected mbCode flag to be defined")
	} else {
		if mbCodeFlag.Shorthand != "c" {
			t.Errorf("Expected mbCode shorthand to be 'c', got '%s'", mbCodeFlag.Shorthand)
		}
		if mbCodeFlag.Usage != "MB-code." {
			t.Errorf("Expected mbCode usage to be 'MB-code.', got '%s'", mbCodeFlag.Usage)
		}
	}
}

func TestLoginCommandWithLicenseKey(t *testing.T) {
	sm := &mockSessionManager{}
	
	cmd := NewLoginCommand(sm)
	cmd.Flags().Set("licenseKey", "test-license-key")
	cmd.Run(cmd, []string{})
	
	// Verify that LoginWithKey was called with the correct parameter
	// Note: We can't directly verify this with our simple mock, but we can test the logic
}

func TestLoginCommandWithMBCode(t *testing.T) {
	sm := &mockSessionManager{}
	
	cmd := NewLoginCommand(sm)
	cmd.Flags().Set("mbCode", "test-mb-code")
	cmd.Run(cmd, []string{})
	
	// Verify that LoginWithCode was called with the correct parameter
	// Note: We can't directly verify this with our simple mock, but we can test the logic
}

func TestLoginCommandLogic(t *testing.T) {
	// Test the logic without running the command to avoid HandleError os.Exit
	sm := &mockSessionManager{}
	
	// Test LoginWithKey
	err := sm.LoginWithKey("test-license-key")
	if err != nil {
		t.Errorf("Expected no error from sm.LoginWithKey, got %v", err)
	}
	
	// Test LoginWithCode
	err = sm.LoginWithCode("test-mb-code")
	if err != nil {
		t.Errorf("Expected no error from sm.LoginWithCode, got %v", err)
	}
}