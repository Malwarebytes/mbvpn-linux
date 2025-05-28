/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/errors"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
	"github.com/Malwarebytes/mbvpn/pkg/servers"
	"github.com/Malwarebytes/mbvpn/pkg/session"
	"github.com/Malwarebytes/mbvpn/pkg/vpn"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

// ErrorHandler is the central error handler for all commands
var ErrorHandler *errors.Handler

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
		trace, _ := cmd.Flags().GetBool("trace")

		// Set log level based on flags
		switch {
		case trace:
			fmt.Println("=== Running in trace mode ===")
			log.SetLevel(log.TraceLevel)
			config.DebugFlag = true // Trace mode implies debug mode
		case debug:
			fmt.Println("=== Running in debug mode ===")
			log.SetLevel(log.DebugLevel)
			config.DebugFlag = true
		default:
			// Set default log level to info
			log.SetLevel(log.InfoLevel)
			config.DebugFlag = false
		}
		
		// Initialize error handler with current debug setting
		ErrorHandler = errors.NewHandler()
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
	holocron := remote.NewDefaultHolocron(config.NewConfigFileMachineIdProvider())
	sm := session.NewDefaultSessionManager(cp, holocron)
	ss := servers.NewDefaultServerStorage()
	vpn := vpn.NewDefaultVpn(cp, holocron, ss)

	rootCmd.AddCommand(NewLoginCommand(sm))
	rootCmd.AddCommand(NewLogoutCommand(sm))
	rootCmd.AddCommand(NewServersCommand(sm, vpn))
	rootCmd.AddCommand(NewCountriesCommand(sm, vpn))
	rootCmd.AddCommand(NewCitiesCommand(sm, vpn))
	rootCmd.AddCommand(NewConnectCommand(sm, vpn))
	rootCmd.AddCommand(NewDisconnectCommand(vpn))
	rootCmd.AddCommand(NewStatusCommand(vpn))
	rootCmd.AddCommand(NewVersionCommand())

	rootCmd.PersistentFlags().Bool("debug", false, "Run command in debug mode.")
	rootCmd.PersistentFlags().Bool("trace", false, "Run command in trace mode with full request/response logging.")
}

// HandleError processes errors according to their type and debug mode
func HandleError(err error) {
	if err == nil {
		return
	}
	
	if ErrorHandler == nil {
		// Fallback if handler isn't initialized
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	if ErrorHandler.Handle(err) {
		os.Exit(1)
	}
}
