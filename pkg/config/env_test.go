package config

import (
	"testing"
)

func TestDebug(t *testing.T) {
	// Save original flag value and restore it after test
	originalDebugFlag := DebugFlag
	defer func() { DebugFlag = originalDebugFlag }()

	// Test case 1: Debug flag is true
	DebugFlag = true
	if !Debug() {
		t.Error("Debug() should return true when DebugFlag is true")
	}

	// Test case 2: Debug flag is false
	// Note: we're in testing mode, so Debug() should still return true
	DebugFlag = false
	if !Debug() {
		t.Error("Debug() should return true when in testing mode")
	}
}

func TestVerbose(t *testing.T) {
	// Verbose() just calls testing.Testing() and testing.Verbose()
	// Since we're in testing mode, the first condition is true
	// The second condition depends on whether the test was run with -v flag
	// We can't easily mock testing.Verbose(), so we'll just check that the
	// function returns a value without panicking
	_ = Verbose()
}

func TestVersion(t *testing.T) {
	// Save original version values and restore them after test
	origMajor, origMinor, origPatch := VersionMajor, VersionMinor, VersionPatch
	defer func() {
		VersionMajor, VersionMinor, VersionPatch = origMajor, origMinor, origPatch
	}()

	// Test case 1: Default version values
	VersionMajor, VersionMinor, VersionPatch = "1", "2", "3"
	expected := "1.2.3"
	if Version() != expected {
		t.Errorf("Version() = %s, want %s", Version(), expected)
	}

	// Test case 2: Different version values
	VersionMajor, VersionMinor, VersionPatch = "5", "6", "7"
	expected = "5.6.7"
	if Version() != expected {
		t.Errorf("Version() = %s, want %s", Version(), expected)
	}

	// Test case 3: Zero values
	VersionMajor, VersionMinor, VersionPatch = "0", "0", "0"
	expected = "0.0.0"
	if Version() != expected {
		t.Errorf("Version() = %s, want %s", Version(), expected)
	}
}