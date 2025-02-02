package remote

import (
  "testing"
	"github.com/Malwarebytes/mbvpn/config"
)

const(
  licenseKey = "CBKGF-JRFZJ-FZYF8-4QR8E"
)

func TestRegisterDevice(t *testing.T) {
	config.Debug = false
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


	config.Debug = true
  m, err := ActivateDevice(token, licenseKey, false)
	config.Debug = false
  defer DeactivateDevice(token)

  if err != nil {
		t.Fatalf(`ActivateDevice(token, licenseKey, false) = %v`, err)
  }

	if m.Status != "licensed" {
		t.Fatalf(`ActivateDevice(token, licenseKey, false) = { DeviceModule.Status = "%s" }, want "licensed"`, m.Status)
	}
}

func TestDeactivateDevice(t *testing.T) {
	config.Debug = false
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
