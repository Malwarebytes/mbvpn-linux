package cmd

import (
	"fmt"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display the application version.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("v%s\n", config.Version())
		},
	}
}
