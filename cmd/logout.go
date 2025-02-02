/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/session"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout this device.",
	Long: `Deactivates this device seat and cleans local configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
    config.Debug, _ = cmd.Flags().GetBool("debug")
		session.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
