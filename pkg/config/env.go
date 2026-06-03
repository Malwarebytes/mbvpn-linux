package config

import (
	"os"
	"testing"
)

var (
	DebugFlag   bool
	BuildEnv    = "production"
	HolocronUrl = "https://holocron.mwbsys.com/graphql"

	// Version is injected at build time via -ldflags by GoReleaser.
	Version = "dev"
)

func Debug() bool {
	return DebugFlag || testing.Testing()
}

func Verbose() bool {
	return testing.Testing() && testing.Verbose()
}

// GetHolocronUrl returns the Holocron URL, checking environment variable first
func GetHolocronUrl() string {
	if url := os.Getenv("MBVPN_HOLOCRON_URL"); url != "" {
		return url
	}
	return HolocronUrl
}
