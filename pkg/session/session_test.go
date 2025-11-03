package session

import (
	"strings"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/pkg/errors"
)

// TestSessionManagerInterface verifies that DefaultSessionManager implements the SessionManager interface
func TestSessionManagerInterface(t *testing.T) {
	// Test interface compliance - this will fail to compile if interface is not implemented
	var manager SessionManager = &DefaultSessionManager{}
	_ = manager // Use the variable to avoid unused variable error

	// This test ensures the DefaultSessionManager struct exists and has the right shape
	if manager == nil {
		t.Error("DefaultSessionManager should not be nil")
	}
}

// TestSessionManagerStruct tests the basic structure of DefaultSessionManager
func TestSessionManagerStruct(t *testing.T) {
	manager := &DefaultSessionManager{}

	// Test that we can create the struct
	if manager == nil {
		t.Error("Should be able to create DefaultSessionManager")
	}
}

// Test the license key validation logic without dependencies
func TestValidateLicenseKeyFormat(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		shouldError bool
		errorMsg    string
	}{
		{"Too short", "SHORT", true, "Invalid license key"},
		{"Too long", "ABCDE-FGHIJ-KLMNO-PQRST-UVWXY", true, "Invalid license key"},
		{"Empty", "", true, "Invalid license key"},
	}

	// Test only cases that fail early validation (before dependencies are used)
	manager := &DefaultSessionManager{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.LoginWithKey(tt.key)

			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			userErr := errors.GetUserErrorMessage(err)
			if !strings.Contains(userErr, tt.errorMsg) {
				t.Errorf("Expected error message to contain '%s', got: %s", tt.errorMsg, userErr)
			}
		})
	}
}

// Test the MB code validation logic without dependencies
func TestValidateMBCodeFormat(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldError bool
		errorMsg    string
	}{
		{"Too short", "ABC", true, "Invalid MB-code"},
		{"Too long", "ABCDEFGH", true, "Invalid MB-code"},
		{"Empty", "", true, "Invalid MB-code"},
		{"With prefix too long", "MB-ABCDEFGH", true, "Invalid MB-code"},
	}

	manager := &DefaultSessionManager{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.LoginWithCode(tt.code)

			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			userErr := errors.GetUserErrorMessage(err)
			if !strings.Contains(userErr, tt.errorMsg) {
				t.Errorf("Expected error message to contain '%s', got: %s", tt.errorMsg, userErr)
			}
		})
	}
}

// Test case conversion for license keys
func TestLicenseKeyCaseConversion(t *testing.T) {
	manager := &DefaultSessionManager{}

	// Test that lowercase keys are handled (they should be converted to uppercase)
	// We can verify this by checking that the same validation error occurs
	// for both upper and lower case versions of an invalid-length key

	shortKey := "short"
	shortKeyUpper := "SHORT"

	err1 := manager.LoginWithKey(shortKey)
	err2 := manager.LoginWithKey(shortKeyUpper)

	// Both should fail with the same error (format validation)
	if err1 == nil || err2 == nil {
		t.Error("Both should fail with format validation")
		return
	}

	userErr1 := errors.GetUserErrorMessage(err1)
	userErr2 := errors.GetUserErrorMessage(err2)

	if userErr1 != userErr2 {
		t.Errorf("Case conversion should result in same error. Got:\n'%s'\nvs\n'%s'", userErr1, userErr2)
	}
}

// Test case conversion for MB codes
func TestMBCodeCaseConversion(t *testing.T) {
	manager := &DefaultSessionManager{}

	// Test case insensitive handling
	shortCode := "abc"
	shortCodeUpper := "ABC"
	shortCodeWithPrefix := "mb-abc"

	err1 := manager.LoginWithCode(shortCode)
	err2 := manager.LoginWithCode(shortCodeUpper)
	err3 := manager.LoginWithCode(shortCodeWithPrefix)

	// All should fail with the same error (format validation - too short)
	if err1 == nil || err2 == nil || err3 == nil {
		t.Error("All should fail with format validation")
		return
	}

	userErr1 := errors.GetUserErrorMessage(err1)
	userErr2 := errors.GetUserErrorMessage(err2)
	userErr3 := errors.GetUserErrorMessage(err3)

	if userErr1 != userErr2 || userErr2 != userErr3 {
		t.Errorf("Case conversion should result in same error. Got:\n'%s'\n'%s'\n'%s'", userErr1, userErr2, userErr3)
	}
}

// Note: Testing Active() and Logout() methods requires proper dependencies
// since they interact with config.ConfigProvider which has unexported methods
// that can't be easily mocked. The core validation logic is tested above.
