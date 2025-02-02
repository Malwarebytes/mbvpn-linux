package session

import (
	"fmt"
	"log"

	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/remote"
)

func Login(key string, mbcode bool) {
	// Check current session
	if Active() {
		fmt.Println("There is an active session on your device. Try logout command first if you want to re-login.")
		return
	}

	fmt.Println("Welcome to Malwarebytes VPN client!")

	// Register device
	installationToken, err := remote.RegisterDevice()
	if err != nil {
		log.Panic(err)
	}
	config.StoreInstallationToken(installationToken)

	// Activate device
	m, err := remote.ActivateDevice(installationToken, key, mbcode)
	if err != nil {
		println("Cannot activate this device.")
		config.DeleteConfig()
		log.Panic(err)
	}
	fmt.Printf("License status: %s\n", m.Status)
}

func Logout() {
	fmt.Println("Logging out...")

	installationToken, err := config.GetInstallationToken()
	if err == nil {
		m, err := remote.DeactivateDevice(installationToken)
		if err != nil {
			println("Something goes wrong with deactivation. Visit my.malwarebytes.com.")
		}
		fmt.Printf("License status: %s\n", m.Status)
	}

	err = config.DeleteConfig()
	if err != nil {
		println("There is no active session on your device.")
	}
}

func Active() bool {
	installationToken, _ := config.GetInstallationToken()
	return installationToken != ""
}
