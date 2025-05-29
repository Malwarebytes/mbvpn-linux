package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Malwarebytes/mbvpn/pkg/config"
)

func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()
	
	if cmd.Use != "version" {
		t.Errorf("Expected Use to be 'version', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "Display version information" {
		t.Errorf("Expected Short to be 'Display version information', got '%s'", cmd.Short)
	}
	
	expectedLong := "Display the application version, build environment and build type."
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, cmd.Long)
	}
	
	if cmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestVersionCommandOutput(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Run the command
	cmd := NewVersionCommand()
	cmd.Run(cmd, []string{})
	
	// Restore stdout and get output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	
	// Verify output contains expected fields
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines of output, got %d", len(lines))
	}
	
	expectedPrefixes := []string{
		"Version: ",
		"Environment: ",
		"Build Type: ",
		"Debug Mode: ",
	}
	
	for i, prefix := range expectedPrefixes {
		if i >= len(lines) || !strings.HasPrefix(lines[i], prefix) {
			t.Errorf("Line %d should start with '%s', got '%s'", i+1, prefix, lines[i])
		}
	}
	
	// Verify version format (should be in format "0.0.1+0")
	versionLine := lines[0]
	versionValue := strings.TrimPrefix(versionLine, "Version: ")
	expectedVersion := config.Version()
	if versionValue != expectedVersion {
		t.Errorf("Expected version '%s', got '%s'", expectedVersion, versionValue)
	}
	
	// Verify environment value
	envLine := lines[1]
	envValue := strings.TrimPrefix(envLine, "Environment: ")
	if envValue != config.BuildEnv {
		t.Errorf("Expected environment '%s', got '%s'", config.BuildEnv, envValue)
	}
	
	// Verify build type value
	buildTypeLine := lines[2]
	buildTypeValue := strings.TrimPrefix(buildTypeLine, "Build Type: ")
	if buildTypeValue != config.BuildType {
		t.Errorf("Expected build type '%s', got '%s'", config.BuildType, buildTypeValue)
	}
	
	// Verify debug mode value
	debugLine := lines[3]
	debugValue := strings.TrimPrefix(debugLine, "Debug Mode: ")
	expectedDebug := "true" // Since we're in testing mode, Debug() returns true
	if debugValue != expectedDebug {
		t.Errorf("Expected debug mode '%s', got '%s'", expectedDebug, debugValue)
	}
}