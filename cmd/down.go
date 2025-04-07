/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
)

func NewDownCommand(vpn vpn.Vpn) *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect",
		Aliases: []string{"d"},
		Short: "Stops the VPN connection.",
		Long: `Stops the VPN connection.`,
		Run: func(cmd *cobra.Command, args []string) {
			// server := ""
			// if len(args) > 0 {
				// server = args[0]
			// }
			vpn.Down("")
		},
	}
}
