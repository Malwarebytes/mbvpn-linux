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

func NewConnectCommand(sm session.SessionManager, vpn vpn.Vpn) *cobra.Command {
	return &cobra.Command{
		Use:     "connect",
		Aliases: []string{"c"},
		Short:   "Connect to a VPN server",
		Long: `Establishes a VPN connection to the specified server using WireGuard.
Requires an active session (login first) and a valid server identifier.

Server Specification Options:
  You can specify the server in multiple ways:

  1. Exact server name:     us-newyork-1
  2. City code:             NYC (randomly selects from New York servers)
  3. City name:             "New York" (randomly selects from New York servers)
  4. Country code:          US (randomly selects from US servers)
  5. Country name:          "United States" (randomly selects from US servers)

Examples:
  mbvpn connect us-newyork-1     # Connect to specific New York server #1
  mbvpn connect NYC              # Connect to random New York server
  mbvpn connect "New York"       # Connect to random New York server
  mbvpn connect US               # Connect to random US server
  mbvpn connect "United States"  # Connect to random US server

Use 'mbvpn servers' to see all available servers with their exact names.
Use 'mbvpn countries' or 'mbvpn cities' to browse servers by location.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if sm.Active() {
				err := vpn.Connect(args[0])
				HandleError(err)
			} else {
				err := errors.NewUserError("There is no active session on your device. Try 'login' command first.", errors.ErrUnauthorized)
				HandleError(err)
			}
		},
	}
}
