/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/remote"
	"github.com/Malwarebytes/mbvpn/session"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout this device.",
	Long:  `Deactivates this device seat and cleans local configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		sm := session.NewDefaultSessionManager(
			config.NewYamlConfigProvider(),
			remote.NewDefaultHolocron(
				config.NewEtcFileMachineIdProvider(),
			),
		)

		sm.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
