package config

import "testing"

var (
	DebugFlag   bool
	BuildEnv    string
	BuildType   string
	HolocronUrl string
)

func Debug() bool {
	return DebugFlag || testing.Testing()
}

func Verbose() bool {
	return testing.Testing() && testing.Verbose()
}
