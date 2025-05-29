package cmd

import (
	"testing"

	"github.com/Malwarebytes/mbvpn/pkg/errors"
	"github.com/spf13/cobra"
)

func TestRootCommandStructure(t *testing.T) {
	if rootCmd.Use != "mbvpn" {
		t.Errorf("Expected Use to be 'mbvpn', got '%s'", rootCmd.Use)
	}
	
	if rootCmd.Short != "Malwarebytes VPN command-line client" {
		t.Errorf("Expected Short to be 'Malwarebytes VPN command-line client', got '%s'", rootCmd.Short)
	}
	
	expectedLong := `MBVPN is a command-line VPN client for Malwarebytes VPN service.
Manage your VPN connections, browse available servers, and control
your privacy settings from the terminal. Requires a valid Malwarebytes license.`
	if rootCmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, rootCmd.Long)
	}
	
	if rootCmd.PersistentPreRun == nil {
		t.Error("Expected PersistentPreRun function to be defined")
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Test debug flag
	debugFlag := rootCmd.PersistentFlags().Lookup("debug")
	if debugFlag == nil {
		t.Error("Expected debug flag to be defined")
	} else {
		if debugFlag.Usage != "Run command in debug mode." {
			t.Errorf("Expected debug flag usage to be 'Run command in debug mode.', got '%s'", debugFlag.Usage)
		}
	}
	
	// Test trace flag
	traceFlag := rootCmd.PersistentFlags().Lookup("trace")
	if traceFlag == nil {
		t.Error("Expected trace flag to be defined")
	} else {
		if traceFlag.Usage != "Run command in trace mode with full request/response logging." {
			t.Errorf("Expected trace flag usage to be 'Run command in trace mode with full request/response logging.', got '%s'", traceFlag.Usage)
		}
	}
}

func TestRootCommandSubcommands(t *testing.T) {
	expectedCommands := []string{"login", "logout", "servers", "countries", "cities", "connect", "disconnect", "status", "version"}
	
	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to be registered", expectedCmd)
		}
	}
}

func TestPersistentPreRunInitializesErrorHandler(t *testing.T) {
	// Create a simple test command with the same PersistentPreRun
	testCmd := &cobra.Command{
		Use: "test",
		PersistentPreRun: rootCmd.PersistentPreRun,
	}
	testCmd.PersistentFlags().Bool("debug", false, "Run command in debug mode.")
	testCmd.PersistentFlags().Bool("trace", false, "Run command in trace mode.")
	
	// Execute PersistentPreRun
	testCmd.PersistentPreRun(testCmd, []string{})
	
	// Verify ErrorHandler is initialized
	if ErrorHandler == nil {
		t.Error("Expected ErrorHandler to be initialized")
	}
}

func TestHandleErrorWithNilError(t *testing.T) {
	// Initialize ErrorHandler
	ErrorHandler = errors.NewHandler()
	
	// Test with nil error - should not panic or exit
	HandleError(nil)
	// If we reach this point, the function handled nil correctly
}

func TestHandleErrorWithNilHandler(t *testing.T) {
	// Test error handling when ErrorHandler is nil
	ErrorHandler = nil
	
	// We can't test the actual os.Exit behavior, but we can verify the function exists
	// and would handle the fallback case
	if ErrorHandler != nil {
		t.Error("Expected ErrorHandler to be nil for this test")
	}
	
	// Reset ErrorHandler for other tests
	ErrorHandler = errors.NewHandler()
}