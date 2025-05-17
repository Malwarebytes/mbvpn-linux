//go:build integration

package main

import (
	"os"
)

const (
	machineId = "JDfhsjdfKDJfhgdj8fsf3"
)

// getLicenseKey returns the license key from environment variable
func getLicenseKey() string {
	licenseKey := os.Getenv("MBVPN_TEST_LICENSE_KEY")
	if licenseKey == "" {
		panic("MBVPN_TEST_LICENSE_KEY environment variable is required but not set")
	}
	return licenseKey
}
