package config

import (
	"fmt"
	"testing"
)

var (
	DebugFlag   bool
	BuildEnv    string
	HolocronUrl string

	// Version information
	VersionMajor = "0"
	VersionMinor = "0"
	VersionPatch = "1"
	VersionBuild = "0"
)

func Debug() bool {
	return DebugFlag || testing.Testing()
}

func Verbose() bool {
	return testing.Testing() && testing.Verbose()
}

// Version returns the full version string in the format "major.minor.patch+build"
func Version() string {
	return fmt.Sprintf("%s.%s.%s+%s", VersionMajor, VersionMinor, VersionPatch, VersionBuild)
}
