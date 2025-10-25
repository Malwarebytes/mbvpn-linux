/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Malwarebytes/mbvpn-linux/pkg/config"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display the application version, build environment and build type.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", config.Version())
			fmt.Printf("Environment: %s\n", config.BuildEnv)
			fmt.Printf("Build Type: %s\n", config.BuildType)
			fmt.Printf("Debug Mode: %t\n", config.Debug())
		},
	}
}
