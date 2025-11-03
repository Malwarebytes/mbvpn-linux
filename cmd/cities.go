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

func NewCitiesCommand(sm session.SessionManager, vpn vpn.Vpn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cities",
		Short: "List available cities with VPN servers",
		Long: `Displays a list of cities where VPN servers are located.
Provides more specific geographic options than countries. Requires an active session.`,
		Run: func(cmd *cobra.Command, args []string) {
			if sm.Active() {
				err := vpn.Servers(true, false)
				HandleError(err)
			} else {
				err := errors.NewUserError("There is no active session on your device. Try 'login' command first.", errors.ErrUnauthorized)
				HandleError(err)
			}
		},
	}

	return cmd
}
