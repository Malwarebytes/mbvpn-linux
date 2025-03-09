//go:build integration

package main

import (
  "testing"

  "github.com/Malwarebytes/mbvpn/pkg/config"
)

func TestGet(t *testing.T) {
  provider := config.NewEtcFileMachineIdProvider()
  id, err := provider.Get()
  if err != nil {
    t.Fatal(err)
  }

  if id == "" {
    t.Fatal(`Cannot read "/etc/machine-id"`)
  }
} 
