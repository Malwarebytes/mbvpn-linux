package session

import (
	"fmt"
	"log"
	"strings"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/remote"
)

type SessionManager interface {
	LoginWithKey(string)
	LoginWithCode(string)
	login(string, bool)
	Logout()
	Active() bool
}

type DefaultSessionManager struct {
	cfgProvider config.ConfigProvider
	holocron    remote.Holocron
}

func NewDefaultSessionManager(cfgProvider config.ConfigProvider, holocron remote.Holocron) SessionManager {
	return &DefaultSessionManager{
		cfgProvider: cfgProvider,
		holocron:    holocron,
	}
}

func (sm *DefaultSessionManager) LoginWithKey(key string) {
	formattedKey := strings.ToUpper(key)

	if len(formattedKey) != 23 {
		fmt.Println("Invalid license key. License key should be 23 characters long.")
		return
	}

	sm.login(formattedKey, false)
}

func (sm *DefaultSessionManager) LoginWithCode(code string) {
	formattedCode := strings.ToUpper(code)
	formattedCode = strings.TrimPrefix(formattedCode, "MB-")

	if len(formattedCode) != 6 {
		fmt.Println("Invalid MB-code. MB-code should look like MB-XXXXXX or XXXXXX.")
		return
	}

	sm.login(formattedCode, true)
}

func (sm *DefaultSessionManager) login(token string, mbcode bool) {
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
	err = sm.cfgProvider.StoreInstallationToken(installationToken)
	if err != nil {
		log.Panic(err)
		return
	}

	// Activate device
	m, err := sm.holocron.ActivateDevice(installationToken, token, mbcode)
	if err != nil {
		fmt.Println("Cannot activate this device.")
		sm.cfgProvider.DeleteConfig()
		log.Panic(err)
	}
	if m.Status != remote.DeviceStatusLicensed && m.Status != remote.DeviceStatusTrial {
		fmt.Println("Cannot activate this device. Check your license.")
		sm.cfgProvider.DeleteConfig()
		if config.Debug() {
			log.Panic(fmt.Errorf("device status: %s", m.Status))
		}
		return
	}
	fmt.Printf("License status: %s\n", m.Status)
}

func (sm *DefaultSessionManager) Logout() {
	fmt.Println("Logging out...")

	installationToken, err := sm.cfgProvider.GetInstallationToken()
	if err == nil {
		sm.holocron.DeactivateDevice(installationToken)
	}

	err = sm.cfgProvider.DeleteConfig()
	if err != nil {
		fmt.Println("There is no active session on your device.")
	}
}

func (sm *DefaultSessionManager) Active() bool {
	installationToken, err := sm.cfgProvider.GetInstallationToken()
	return installationToken != "" && err == nil
}
