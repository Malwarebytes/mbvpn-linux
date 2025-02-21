package config

import (
	"os"
	"strings"
)

type MachineIdProvider interface {
	Get() (string, error)
}

type EtcFileMachineIdProvider struct{}

func NewEtcFileMachineIdProvider() MachineIdProvider {
	return &EtcFileMachineIdProvider{}
}

func (cp *EtcFileMachineIdProvider) Get() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}

	id := string(data)
	id = strings.TrimSuffix(id, "\n")

	return id, nil
}
