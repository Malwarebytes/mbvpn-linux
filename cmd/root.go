/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
	"github.com/Malwarebytes/mbvpn/pkg/session"
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mbvpn",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			fmt.Println("=== Running in debug mode ===")
		}
		config.DebugFlag = debug
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cp := config.NewYamlConfigProvider()
	holocron := remote.NewDefaultHolocron(config.NewEtcFileMachineIdProvider())
	sm := session.NewDefaultSessionManager(cp, holocron)
	vpn := vpn.NewDefaultVpn(cp, holocron)

	rootCmd.AddCommand(NewLoginCommand(sm))
  rootCmd.AddCommand(NewLogoutCommand(sm))
	rootCmd.AddCommand(NewServersCommand(sm, vpn))
	rootCmd.AddCommand(NewCountriesCommand(sm, vpn))
	rootCmd.AddCommand(NewCitiesCommand(sm, vpn))
  rootCmd.AddCommand(NewUpCommand(sm, vpn))
  rootCmd.AddCommand(NewDownCommand(vpn))
  rootCmd.AddCommand(NewStatusCommand(vpn))

	rootCmd.PersistentFlags().Bool("debug", false, "Run command in debug mode.")
}
