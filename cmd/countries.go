/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/Malwarebytes/mbvpn/pkg/session"
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
)

func NewCountriesCommand(sm session.SessionManager, vpn vpn.Vpn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "countries",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if sm.Active() {
				vpn.Servers(false, false)
			} else {
				fmt.Println(`There is no active session on your device. Try "login" command first.`)
			}
		},
	}

	return cmd
}
