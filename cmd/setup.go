package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/malwarebytes/mbvpn-linux/pkg/output"
	"github.com/spf13/cobra"
)

func NewSetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Set up mbvpn with required capabilities",
		Long: `Sets the CAP_NET_ADMIN capability on the mbvpn binary.

This command must be run as root (using sudo) and only needs to be run once
after installation. After setup, all VPN commands can be run without sudo.

Example:
  sudo mbvpn setup`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetup()
		},
	}

	return cmd
}

func runSetup() error {
	// Check if running as root
	if os.Geteuid() != 0 {
		return fmt.Errorf("setup must be run as root (use: sudo mbvpn setup)")
	}

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	output.PrintMsg(fmt.Sprintf("Setting up capabilities for: %s", execPath), output.MsgOutput)

	// Set CAP_NET_ADMIN capability using setcap
	setcapCmd := exec.Command("setcap", "cap_net_admin+ep", execPath)
	setcapCmd.Stdout = os.Stdout
	setcapCmd.Stderr = os.Stderr

	if err := setcapCmd.Run(); err != nil {
		return fmt.Errorf("failed to set capabilities: %w\n\nMake sure 'setcap' is installed (usually in libcap2-bin package)", err)
	}

	output.PrintMsg("Setup complete. You can now run mbvpn commands without sudo.", output.MsgSuccess)
	output.PrintMsg("", output.MsgOutput)
	output.PrintMsg("Note: If you reinstall or update mbvpn, you'll need to run 'sudo mbvpn setup' again.", output.MsgOutput)

	return nil
}

func init() {
	rootCmd.AddCommand(NewSetupCommand())
}
