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
	origVersion := config.Version
	config.Version = "1.2.3"
	defer func() { config.Version = origVersion }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewVersionCommand()
	cmd.Run(cmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	if strings.Contains(output, "\n") {
		t.Errorf("Expected single line output, got multiple lines: %s", output)
	}

	if !strings.HasPrefix(output, "v") {
		t.Errorf("Expected output to start with 'v', got: %s", output)
	}

	expectedVersion := "v" + config.Version
	if output != expectedVersion {
		t.Errorf("Expected version '%s', got '%s'", expectedVersion, output)
	}

	if !strings.Contains(output, ".") {
		t.Errorf("Expected version format 'vX.Y.Z', got: %s", output)
	}
}
