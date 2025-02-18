/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (

	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/remote"
	"github.com/Malwarebytes/mbvpn/vpn"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		holocron := remote.NewDefaultHolocron(config.NewEtcFileMachineIdProvider())
		cp := config.NewYamlConfigProvider()
		vpn := vpn.NewDefaultVpn(cp, holocron)

		vpn.Down(args[0])
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
