/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/session"
	"github.com/Malwarebytes/mbvpn/vpn"
	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
    config.Debug, _ = cmd.Flags().GetBool("debug")

		if session.Active() {
			vpn.Servers()
		} else {
      fmt.Println(`There is no active session on your device. Try "login" command first.`)
    }
	},
}

func init() {
	rootCmd.AddCommand(serversCmd)
}
