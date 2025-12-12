package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
)

func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd.Use != "version" {
		t.Errorf("Expected Use to be 'version', got '%s'", cmd.Use)
	}

	if cmd.Short != "Display version information" {
		t.Errorf("Expected Short to be 'Display version information', got '%s'", cmd.Short)
	}

	expectedLong := "Display the application version."
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
	output := strings.TrimSpace(buf.String())

	// Verify output is single line
	if strings.Contains(output, "\n") {
		t.Errorf("Expected single line output, got multiple lines: %s", output)
	}

	// Verify output starts with "v"
	if !strings.HasPrefix(output, "v") {
		t.Errorf("Expected output to start with 'v', got: %s", output)
	}

	// Verify version format matches "vX.Y.Z"
	expectedVersion := "v" + config.Version()
	if output != expectedVersion {
		t.Errorf("Expected version '%s', got '%s'", expectedVersion, output)
	}

	// Verify it contains dots (semantic version format)
	if !strings.Contains(output, ".") {
		t.Errorf("Expected version format 'vX.Y.Z', got: %s", output)
	}
}
