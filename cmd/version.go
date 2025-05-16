package cmd

import (
	"fmt"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command
func NewVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display detailed version information about the MBVPN build`,
		Run: func(cmd *cobra.Command, args []string) {
			showVersion()
		},
	}

	return versionCmd
}

func showVersion() {
	fmt.Println("MBVPN Version Information")
	fmt.Println("------------------------")
	fmt.Printf("Environment: %s\n", getEnv())
	fmt.Printf("Build Type:  %s\n", getBuildType())
	fmt.Printf("Debug Mode:  %t\n", config.Debug())
}

func getEnv() string {
	if config.BuildEnv == "" {
		return "development"
	}
	return config.BuildEnv
}

func getBuildType() string {
	if config.BuildType == "" {
		return "debug"
	}
	return config.BuildType
}