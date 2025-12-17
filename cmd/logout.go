package cmd

import (
	"github.com/malwarebytes/mbvpn-linux/pkg/session"
	"github.com/spf13/cobra"
)

func NewLogoutCommand(sm session.SessionManager) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout this device.",
		Long:  `Deactivates this device seat and cleans local configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := sm.Logout()
			HandleError(err)
		},
	}
}
