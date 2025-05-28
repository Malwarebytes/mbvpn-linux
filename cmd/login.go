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
		Long: `Uses provided credentials to perform activation for this device.
   The command claims available seat of your Malwarebytes license.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			
			key, _ := cmd.Flags().GetString("licenseKey")
			if key == "" {
				code, _ := cmd.Flags().GetString("mbCode")
				err = sm.LoginWithCode(code)
			} else {
				err = sm.LoginWithKey(key)
			}
			
			// Handle any errors that might have occurred
			HandleError(err)
		},
	}

	cmd.Flags().StringP("licenseKey", "k", "", "License key.")
	cmd.Flags().StringP("mbCode", "c", "", "MB-code.")
	cmd.MarkFlagsOneRequired("licenseKey", "mbCode")

	return cmd
}