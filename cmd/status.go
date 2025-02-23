/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the current VPN connection",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		holocron := remote.NewDefaultHolocron(config.NewEtcFileMachineIdProvider())
		cp := config.NewYamlConfigProvider()
		vpn := vpn.NewDefaultVpn(cp, holocron)

		vpn.Status()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
