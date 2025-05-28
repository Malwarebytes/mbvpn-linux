/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/pkg/errors"
	"github.com/Malwarebytes/mbvpn/pkg/session"
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
)

func NewCountriesCommand(sm session.SessionManager, vpn vpn.Vpn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "countries",
		Short: "List available countries with VPN servers",
		Long: `Displays a list of countries where VPN servers are available.
Useful for selecting servers by geographic location. Requires an active session.`,
		Run: func(cmd *cobra.Command, args []string) {
			if sm.Active() {
				err := vpn.Servers(false, false)
				HandleError(err)
			} else {
				err := errors.NewUserError("There is no active session on your device. Try 'login' command first.", errors.ErrUnauthorized)
				HandleError(err)
			}
		},
	}

	return cmd
}