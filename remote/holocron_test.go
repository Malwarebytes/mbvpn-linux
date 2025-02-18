package remote

import (
	"testing"

	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const(
  licenseKey = "CBKGF-JRFZJ-FZYF8-4QR8E"
  machineId = "JDfhsjdfKDJfhgdj8fsf3"
)

type MockMachineIdProvider struct {
  mId string
}
func (mip *MockMachineIdProvider) Get() (string, error) {
  if mip.mId == "" {
    mip.mId = uuid.New().String()
  }
  return mip.mId, nil
}

type MockConfigProvider struct {
  it string
}

func (cp *MockConfigProvider) StoreInstallationToken(token string) {
  cp.it = token
}

func (cp *MockConfigProvider) GetInstallationToken() (string, error) {
  return cp.it, nil
}

func (cp *MockConfigProvider) DeleteConfig() error {
  cp.it = ""
  return nil
}

var h = NewDefaultHolocron(&MockMachineIdProvider{})

func TestRegisterDevice(t *testing.T) {
  token, err := h.RegisterDevice()
  if err != nil {
		t.Fatalf(`RegisterDevice() = %v`, err)
  }
  if token == "" {
		t.Fatalf(`RegisterDevice() = "", want non-empty`)
  }
}

//TODO add test with MBCode
func TestActivateDeviceWithLicenseKey(t *testing.T) {
  installationToken, err := h.RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  m, err := h.ActivateDevice(installationToken, licenseKey, false)
  defer h.DeactivateDevice(installationToken)

  if err != nil {
		t.Fatalf(`ActivateDevice(token, licenseKey, false) = %v`, err)
  }

	if m.Status != "licensed" {
		t.Fatalf(`ActivateDevice(token, licenseKey, false) = { DeviceModule.Status = "%s" }, want "licensed"`, m.Status)
	}
}

func TestDeactivateDevice(t *testing.T) {
  installationToken, err := h.RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  _, err = h.ActivateDevice(installationToken, licenseKey, false)
	if err != nil {
		t.Fatal(err)
	}

  m, err := h.DeactivateDevice(installationToken)
  if err != nil {
		t.Fatalf(`DeactivateDevice(token) = %v`, err)
  }
	if m.Status != "unknown" {
		t.Fatalf(`DeactivateDevice(token) = { DeviceModule.Status = "%s" }, want "unknown"`, m.Status)
	}
}

func TestVpnRegisterPublicKey(t *testing.T) {
  installationToken, err := h.RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  _, err = h.ActivateDevice(installationToken, licenseKey, false)
  if err != nil {
    t.Fatal(err)
  }
  defer h.DeactivateDevice(installationToken)

  wgKeys, err := wgtypes.GeneratePrivateKey()
  if err != nil {
    t.Fatal(err)
  }
  wgPubKey := wgKeys.PublicKey().String()

  ip, err := h.VpnRegisterPublicKey(installationToken, wgPubKey)
  if err != nil {
    t.Fatalf(`VpnRegisterPublicKey(token, "%s") = %v`, wgPubKey, err)
  }
  if ip.IpV4 == "" || ip.IpV6 == "" || ip.KeyExpirationHours == 0 {
    t.Fatalf(`VpnRegisterPublicKey(token, "%s") = { IpV6 = %s, IpV6 = %s, KeyExpirationHours = %d }, want non-empty IPs and KeyExpirationHours > 0`, wgPubKey, ip.IpV4, ip.IpV6, ip.KeyExpirationHours)
  }
}

func TestGetVpnNetworkDetails(t *testing.T) {
  _, err := h.RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  d, err := h.GetVpnNetworkDetails()
  if err != nil {
		t.Fatalf(`GetVpnNetworkDetails(token) = %v`, err)
  }
	if d.Ip == "" {
		t.Fatal(`GetVpnNetworkDetails(token) = { ip = <EMPTY> }, want to be present`)
	}
}

func TestGetVpnLocations(t *testing.T) {
  _, err := h.RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  l, err := h.GetVpnLocations()
  if err != nil {
    t.Fatalf(`GetVpnLocations(token) = %v`, err)
  }

  if len(l.Countries) == 0 {
    t.Fatal(`GetVpnLocations(token) = no countries, want countries array`)
  }
}
