package cmd

import (
	"github.com/malwarebytes/mbvpn-linux/pkg/session"
	"github.com/spf13/cobra"
)

func NewLoginCommand(sm session.SessionManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Activates this device with your Malwarebytes license.",
		Long: `Uses provided credentials to perform activation for this device.
   The command claims available seat of your Malwarebytes license.`,
		Run: func(cmd *cobra.Command, args []string) {
			code, _ := cmd.Flags().GetString("code")
			err := sm.LoginWithCode(code)
			HandleError(err)
		},
	}

	cmd.Flags().StringP("code", "c", "", "MB-code.")
	cmd.MarkFlagsOneRequired("code")

	return cmd
}
