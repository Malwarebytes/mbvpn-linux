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

func NewServersCommand(sm session.SessionManager, vpn vpn.Vpn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servers",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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