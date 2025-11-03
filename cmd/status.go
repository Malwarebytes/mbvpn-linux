/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/malwarebytes/mbvpn-linux/pkg/vpn"
	"github.com/spf13/cobra"
)

func NewStatusCommand(vpn vpn.Vpn) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of the current VPN connection",
		Long: `Displays detailed information about the current VPN connection status,
including connection state, server details, and network configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := vpn.Status()
			HandleError(err)
		},
	}
}
