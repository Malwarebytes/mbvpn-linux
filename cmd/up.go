/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/remote"
	"github.com/Malwarebytes/mbvpn/session"
	"github.com/Malwarebytes/mbvpn/vpn"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		holocron := remote.NewDefaultHolocron(config.NewEtcFileMachineIdProvider())
    cp := config.NewYamlConfigProvider()
		sm := session.NewDefaultSessionManager(cp, holocron)
		vpn := vpn.NewDefaultVpn(cp, holocron)

		if sm.Active() {
			vpn.Up(args[0])
		} else {
			fmt.Println(`There is no active session on your device. Try "login" command first.`)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
