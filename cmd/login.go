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
			key, _ := cmd.Flags().GetString("licenseKey")
			if key == "" {
				code, _ := cmd.Flags().GetString("mbCode")
				sm.LoginWithCode(code)
			} else {
				sm.LoginWithKey(key)
			}
		},
	}

	cmd.Flags().StringP("licenseKey", "k", "", "License key.")
	cmd.Flags().StringP("mbCode", "c", "", "MB-code.")
	cmd.MarkFlagsOneRequired("licenseKey", "mbCode")

	return cmd
}
