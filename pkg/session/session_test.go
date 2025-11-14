package session

import (
	"strings"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/pkg/errors"
)

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

func TestMBCodeCaseConversion(t *testing.T) {
	manager := &DefaultSessionManager{}

	shortCode := "abc"
	shortCodeUpper := "ABC"
	shortCodeWithPrefix := "mb-abc"

	err1 := manager.LoginWithCode(shortCode)
	err2 := manager.LoginWithCode(shortCodeUpper)
	err3 := manager.LoginWithCode(shortCodeWithPrefix)

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
