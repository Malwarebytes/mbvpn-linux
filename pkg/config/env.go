package config

import "testing"

var DebugFlag bool

func Debug() bool {
	return DebugFlag || testing.Testing()
}

func Verbose() bool {
	return testing.Testing() && testing.Verbose()
}
