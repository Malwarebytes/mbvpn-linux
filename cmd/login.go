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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Activates this device with your Malwarebytes license.",
	Long: `Uses provided creadentials to perform activation for this device.
   The command claimes available seat of your Malwarebytes license.`,
	Run: func(cmd *cobra.Command, args []string) {
		mbcode := false
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			key, _ = cmd.Flags().GetString("mbcode")
			mbcode = true
		}

		sm := session.NewDefaultSessionManager(
			config.NewYamlConfigProvider(),
			remote.NewDefaultHolocron(
				config.NewEtcFileMachineIdProvider(),
			),
		)

		sm.Login(key, mbcode)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("key", "k", "", "License key.")
	loginCmd.Flags().StringP("mbcode", "m", "", "MB-code.")
	loginCmd.MarkFlagsOneRequired("key", "mbcode")
}
