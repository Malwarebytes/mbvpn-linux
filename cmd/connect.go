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
  You can specify the server in multiple ways (case-insensitive):

  1. Exact server name:     gb-lon-2
  2. City code:             gb-lon (randomly selects from London servers)
  3. City name:             "London" (randomly selects from London servers)
  4. Country code:          gb (randomly selects from UK servers)
  5. Country name:          "United Kingdom" (randomly selects from UK servers)

Examples:
  mbvpn connect gb-lon-2         # Connect to specific London server #2
  mbvpn connect gb-lon           # Connect to random London server
  mbvpn connect London           # Connect to random London server
  mbvpn connect gb               # Connect to random UK server
  mbvpn connect "United Kingdom" # Connect to random UK server

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
