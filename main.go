/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/Malwarebytes/mbvpn/cmd"
	"github.com/Malwarebytes/mbvpn/config"
)

var (
  Debug bool
)

func main() {
  config.Debug = Debug
  if Debug {
    println("=== Debug mode ===")
  }
	cmd.Execute()
}
