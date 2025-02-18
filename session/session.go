package session

import (
	"fmt"
	"log"

	"github.com/Malwarebytes/mbvpn/config"
	"github.com/Malwarebytes/mbvpn/remote"
)

type SessionManager interface {
  Login(string, bool)
  Logout()
  Active() bool
}

type DefaultSessionManager struct {
  cfgProvider config.ConfigProvider
  holocron remote.Holocron
}

func NewDefaultSessionManager(cfgProvider config.ConfigProvider, holocron remote.Holocron) SessionManager {
  return &DefaultSessionManager{
    cfgProvider: cfgProvider,
    holocron: holocron,
  }
}

func (sm *DefaultSessionManager) Login(key string, mbcode bool) {
	// Check current session
	if sm.Active() {
		fmt.Println("There is an active session on your device. Try logout command first if you want to re-login.")
		return
	}

	fmt.Println("Welcome to Malwarebytes VPN client!")

	// Register device
	installationToken, err := sm.holocron.RegisterDevice()
	if err != nil {
		log.Panic(err)
	}
	sm.cfgProvider.StoreInstallationToken(installationToken)

	// Activate device
	m, err := sm.holocron.ActivateDevice(installationToken, key, mbcode)
	if err != nil {
		println("Cannot activate this device.")
		sm.cfgProvider.DeleteConfig()
		log.Panic(err)
	}
	fmt.Printf("License status: %s\n", m.Status)
}

func (sm *DefaultSessionManager) Logout() {
	fmt.Println("Logging out...")

	installationToken, err := sm.cfgProvider.GetInstallationToken()
	if err == nil {
		m, err := sm.holocron.DeactivateDevice(installationToken)
		if err != nil {
			println("Something goes wrong with deactivation. Visit my.malwarebytes.com.")
		}
		fmt.Printf("License status: %s\n", m.Status)
	}

	err = sm.cfgProvider.DeleteConfig()
	if err != nil {
		println("There is no active session on your device.")
	}
}

func (sm *DefaultSessionManager) Active() bool {
	installationToken, _ := sm.cfgProvider.GetInstallationToken()
	return installationToken != ""
}
