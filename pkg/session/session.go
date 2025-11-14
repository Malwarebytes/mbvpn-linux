package session

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	"github.com/malwarebytes/mbvpn-linux/pkg/errors"
	"github.com/malwarebytes/mbvpn-linux/pkg/output"
	"github.com/malwarebytes/mbvpn-linux/pkg/remote"
)

type SessionManager interface {
	LoginWithCode(code string) error
	Logout() error
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

func (sm *DefaultSessionManager) LoginWithCode(code string) error {
	formattedCode := strings.ToUpper(code)
	formattedCode = strings.TrimPrefix(formattedCode, "MB-")

	if len(formattedCode) != 6 {
		return errors.NewUserError("Invalid MB-code. MB-code should look like MB-XXXXXX or XXXXXX.", errors.ErrInvalidInput)
	}

	return sm.login(formattedCode)
}

func (sm *DefaultSessionManager) login(code string) error {
	// Check current session
	if sm.Active() {
		return errors.NewUserError("There is an active session on your device. Try logout command first if you want to re-login.", nil)
	}

	output.PrintMsg("Welcome to Malwarebytes VPN client!", output.MsgOutput)

	// Register device
	installationToken, err := sm.holocron.RegisterDevice()
	if err != nil {
		log.Errorf("Error registering the device: %v", err)
		return errors.NewNetworkError("register device", err)
	}

	err = sm.cfgProvider.StoreInstallationToken(installationToken)
	if err != nil {
		log.Errorf("Error storing installation token: %v", err)
		return errors.NewConfigError("store installation token", err)
	}

	// Activate device
	m, err := sm.holocron.ActivateDevice(installationToken, code)
	if err != nil {
		sm.cfgProvider.DeleteConfig()
		log.Errorf("Error activating the device: %v", err)
		return errors.NewNetworkError("activate device", err)
	}

	if m.Status != remote.DeviceStatusLicensed && m.Status != remote.DeviceStatusTrial {
		_ = sm.cfgProvider.DeleteConfig()
		err := fmt.Errorf("device status: %s", m.Status)
		return errors.NewAuthError("check license status", err)
	}

	output.PrintMsg("Activated successfully!", output.MsgSuccess)
	output.PrintMsg(fmt.Sprintf("License status: %s", m.Status), output.MsgSuccess)
	return nil
}

func (sm *DefaultSessionManager) Logout() error {
	installationToken, err := sm.cfgProvider.GetInstallationToken()
	if err == nil {
		_, deactivateErr := sm.holocron.DeactivateDevice(installationToken)
		if deactivateErr != nil {
			log.Errorf("Error deactivating device: %v", deactivateErr)
			// Continue with logout even if deactivation fails
		}
	}

	err = sm.cfgProvider.DeleteConfig()
	if err != nil {
		return errors.NewUserError("There is no active session on your device.", err)
	}

	output.PrintMsg("Logged out successfully!", output.MsgSuccess)
	return nil
}

func (sm *DefaultSessionManager) Active() bool {
	installationToken, err := sm.cfgProvider.GetInstallationToken()
	return installationToken != "" && err == nil
}
