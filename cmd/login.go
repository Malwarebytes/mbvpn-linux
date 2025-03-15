/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Malwarebytes/mbvpn/pkg/session"
	"github.com/spf13/cobra"
)

func NewLoginCommand(sm session.SessionManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Activates this device with your Malwarebytes license.",
		Long: `Uses provided creadentials to perform activation for this device.
   The command claimes available seat of your Malwarebytes license.`,
		Run: func(cmd *cobra.Command, args []string) {
			mbcode := false
			key, _ := cmd.Flags().GetString("key")
			if key == "" {
				key, _ = cmd.Flags().GetString("mbcode")
				mbcode = true
			}

			sm.Login(key, mbcode)
		},
	}

	cmd.Flags().StringP("key", "k", "", "License key.")
	cmd.Flags().StringP("mbcode", "m", "", "MB-code.")
	cmd.MarkFlagsOneRequired("key", "mbcode")

	return cmd
}
