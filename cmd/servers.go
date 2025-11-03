/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/malwarebytes/mbvpn-linux/pkg/errors"
	"github.com/malwarebytes/mbvpn-linux/pkg/session"
	"github.com/malwarebytes/mbvpn-linux/pkg/vpn"
	"github.com/spf13/cobra"
)

func NewServersCommand(sm session.SessionManager, vpn vpn.Vpn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servers",
		Short: "List all available VPN servers",
		Long: `Displays a comprehensive list of all available VPN servers with their
locations, names, and connection details. Requires an active session.`,
		Run: func(cmd *cobra.Command, args []string) {
			if sm.Active() {
				err := vpn.Servers(true, true)
				HandleError(err)
			} else {
				err := errors.NewUserError("There is no active session on your device. Try 'login' command first.", errors.ErrUnauthorized)
				HandleError(err)
			}
		},
	}

	return cmd
}
