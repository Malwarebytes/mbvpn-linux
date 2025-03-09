/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/Malwarebytes/mbvpn/cmd"
)

func main() {
	// systray.Run(onReady, onExit)
  cmd.Execute()
}

// func onReady() {
// 	systray.SetIcon(icon.Data)
// 	systray.SetTitle("MBVPN")
// 	systray.SetTooltip("Malwarebytes VPN")
// 	mQuit := systray.AddMenuItem("Quit", "Quit the app")
//
// 	go func() {
// 		<-mQuit.ClickedCh
// 		systray.Quit()
// 	}()
//
// 	cmd.Execute()
// }
//
// func onExit() {
// }

// Environment to run the systray app
// nix-shell -p libgcc gtk3 libindicator-gtk3 libayatana-appindicator
