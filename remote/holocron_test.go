package remote

import (
	"testing"
)

const(
  licenseKey = "CBKGF-JRFZJ-FZYF8-4QR8E"
  wgPubKey = "38oCmhanYhYMubkpIil8id7ZPqKJsPQX+3JdJJhKE3s="
)

func TestRegisterDevice(t *testing.T) {
  token, err := RegisterDevice()
  if err != nil {
		t.Fatalf(`RegisterDevice() = %v`, err)
  }
  if token == "" {
		t.Fatalf(`RegisterDevice() = "", want non-empty`)
  }
}

//TODO add test with MBCode
func TestActivateDeviceWithLicenseKey(t *testing.T) {
  token, err := RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  m, err := ActivateDevice(token, licenseKey, false)
  defer DeactivateDevice(token)

  if err != nil {
		t.Fatalf(`ActivateDevice(token, licenseKey, false) = %v`, err)
  }

	if m.Status != "licensed" {
		t.Fatalf(`ActivateDevice(token, licenseKey, false) = { DeviceModule.Status = "%s" }, want "licensed"`, m.Status)
	}
}

func TestDeactivateDevice(t *testing.T) {
  token, err := RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  _, err = ActivateDevice(token, licenseKey, false)
	if err != nil {
		t.Fatal(err)
	}

  m, err := DeactivateDevice(token)
  if err != nil {
		t.Fatalf(`DeactivateDevice(token) = %v`, err)
  }
	if m.Status != "unknown" {
		t.Fatalf(`DeactivateDevice(token) = { DeviceModule.Status = "%s" }, want "unknown"`, m.Status)
	}
}

func TestVpnRegisterPublicKey(t *testing.T) {
  token, err := RegisterDevice()
	if err != nil {
		t.Fatal(err)
	}

  _, err = ActivateDevice(token, licenseKey, false)
  if err != nil {
    t.Fatal(err)
  }
  defer DeactivateDevice(token)

  ip, err := VpnRegisterPublicKey(token, wgPubKey)
  if err != nil {
    t.Fatalf(`VpnRegisterPublicKey(token, "%s") = %v`, wgPubKey, err)
  }
  if ip.IpV4 == "" || ip.IpV6 == "" || ip.KeyExpirationHours == 0 {
    t.Fatalf(`VpnRegisterPublicKey(token, "%s") = { IpV6 = %s, IpV6 = %s, KeyExpirationHours = %d }, want non-empty IPs and KeyExpirationHours > 0`, wgPubKey, ip.IpV4, ip.IpV6, ip.KeyExpirationHours)
  }
}
