package session

import (
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"

	"github.com/Malwarebytes/mbvpn/pkg/config"
	"github.com/Malwarebytes/mbvpn/pkg/output"
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
		output.PrintMsg("Invalid license key. License key should be 23 characters long.", output.MsgError)
		return
	}

	sm.login(formattedKey, false)
}

func (sm *DefaultSessionManager) LoginWithCode(code string) {
	formattedCode := strings.ToUpper(code)
	formattedCode = strings.TrimPrefix(formattedCode, "MB-")

	if len(formattedCode) != 6 {
		output.PrintMsg("Invalid MB-code. MB-code should look like MB-XXXXXX or XXXXXX.", output.MsgError)
		return
	}

	sm.login(formattedCode, true)
}

func (sm *DefaultSessionManager) login(token string, mbcode bool) {
	// Check current session
	if sm.Active() {
		output.PrintMsg("There is an active session on your device. Try logout command first if you want to re-login.", output.MsgOutput)
		return
	}

	output.PrintMsg("Welcome to Malwarebytes VPN client!", output.MsgOutput)

	// Register device
	installationToken, err := sm.holocron.RegisterDevice()
	if err != nil {
		log.Errorf("Error registering the device: %v", err)
		return
	}
	err = sm.cfgProvider.StoreInstallationToken(installationToken)
	if err != nil {
		if config.Debug() {
			log.Errorf("Error storing installation token: %v", err)
		}
		return
	}

	// Activate device
	m, err := sm.holocron.ActivateDevice(installationToken, token, mbcode)
	if err != nil {
		output.PrintMsg("Cannot activate this device.", output.MsgError)
		sm.cfgProvider.DeleteConfig()
		log.Errorf("Error activating the device: %v", err)
		return
	}
	if m.Status != remote.DeviceStatusLicensed && m.Status != remote.DeviceStatusTrial {
		output.PrintMsg("Cannot activate this device. Check your license.", output.MsgError)
		_ = sm.cfgProvider.DeleteConfig()
		if config.Debug() {
			log.Panic(fmt.Errorf("device status: %s", m.Status))
		}
		return
	}
	output.PrintMsg("Activated successfully!", output.MsgSuccess)
	output.PrintMsg(fmt.Sprintf("License status: %s", m.Status), output.MsgSuccess)
}

func (sm *DefaultSessionManager) Logout() {
	installationToken, err := sm.cfgProvider.GetInstallationToken()
	if err == nil {
		sm.holocron.DeactivateDevice(installationToken)
	}

	err = sm.cfgProvider.DeleteConfig()
	if err != nil {
		output.PrintMsg("There is no active session on your device.", output.MsgError)
		return
	}

	output.PrintMsg("Logged out successfully!", output.MsgSuccess)
}

func (sm *DefaultSessionManager) Active() bool {
	installationToken, err := sm.cfgProvider.GetInstallationToken()
	return installationToken != "" && err == nil
}
