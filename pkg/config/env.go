package config

import (
	"fmt"
	"os"
	"testing"
)

var (
	DebugFlag   bool
	BuildEnv    = "production"
	HolocronUrl = "https://holocron.mwbsys.com/graphql"

	// Version information
	VersionMajor = "0"
	VersionMinor = "0"
	VersionPatch = "10"
)

func Debug() bool {
	return DebugFlag || testing.Testing()
}

func Verbose() bool {
	return testing.Testing() && testing.Verbose()
}

// Version returns the full version string in the format "major.minor.patch"
func Version() string {
	return fmt.Sprintf("%s.%s.%s", VersionMajor, VersionMinor, VersionPatch)
}

// GetHolocronUrl returns the Holocron URL, checking environment variable first
func GetHolocronUrl() string {
	if url := os.Getenv("MBVPN_HOLOCRON_URL"); url != "" {
		return url
	}
	return HolocronUrl
}
