/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
)

func NewStatusCommand(vpn vpn.Vpn) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of the current VPN connection",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := vpn.Status()
			HandleError(err)
		},
	}
}