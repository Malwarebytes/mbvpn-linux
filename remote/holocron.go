package remote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/Malwarebytes/mbvpn/config"
)

const (
	holocronStageUrl = "https://idms-holocron-stage.mwbsys.com/graphql"
	productCode      = "MBMA-C"
	productVersion   = "5.14.0"
)

var client = &http.Client{Timeout: time.Second * 10}

func RegisterDevice() (string, error) {
	input := RegisterDeviceInput{
		ProductCode:    productCode,
		ProductVersion: productVersion,
		Modules:        []ProductModule{ProductModulePrivacy},
	}
	requestBody := map[string]interface{}{
		"query": `
      mutation RegisterDevice($input: RegisterDeviceInput!) {
        registerDevice(input: $input) {
          device {
            installationToken
          }
        }
      }
    `,
		"variables": map[string]interface{}{
			"input": input,
		},
	}

	var response RegisterDeviceResponse
	err := doRequest("", &requestBody, &response)
	if err != nil {
		return "", err
	}

	return response.Data.RegisterDevice.Device.InstallationToken, nil
}

func ActivateDevice(installationToken string, key string, mbcode bool) (*DeviceModule, error) {
	input := ActivateDeviceInput{
		Modules:          []ProductModule{ProductModulePrivacy},
		ActivationMode:   ActivationModePassive,
	}
	if mbcode {
		input.ActivationMethod = ActivationMethodOneTimeToken
		input.OneTimeToken = key
	} else {
		input.ActivationMethod = ActivationMethodLicenseKey
		input.LicenseKey = key
	}

	requestBody := map[string]interface{}{
		"query": `
      mutation ActivateDevice($input: ActivateDeviceInput!) {
        activateDevice(input: $input) {
          deviceModules {
            privacy {
              status
              termEndsOn
            }
          }
        }
      }
    `,
		"variables": map[string]interface{}{
			"input": input,
		},
	}

	var response ActivateDeviceResponse
	err := doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.ActivateDevice.DeviceModules.Privacy, nil
}

func DeactivateDevice(installationToken string) (*DeviceModule, error) {
	input := DeactivateDeviceInput{
		Modules:          []ProductModule{ProductModulePrivacy},
		DeactivationMode: DeactivationModePassive,
	}

	requestBody := map[string]interface{}{
		"query": `
      mutation DeactivateDevice($input: DeactivateDeviceInput!) {
        deactivateDevice(input: $input) {
          deviceModules {
            privacy {
              status
              termEndsOn
            }
          }
        }
      }
    `,
		"variables": map[string]interface{}{
			"input": input,
		},
	}

	var response DeactivateDeviceResponse
	err := doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.DeactivateDevice.DeviceModules.Privacy, nil
}

func VpnRegisterPublicKey(installationToken string, key string) (*VpnIpAddresses, error) {
	input := VpnRegisterPublicKeyInput{
		PublicKey: key,
	}

	requestBody := map[string]interface{}{
		"query": `
      mutation VpnRegisterPublicKey($input: VpnRegisterPublicKeyInput!) {
        vpnRegisterPublicKey(input: $input) {
          ipv4
          ipv6
          keyExpirationHours
        }
      }
    `,
		"variables": map[string]interface{}{
			"input": input,
		},
	}

	var response VpnRegisterPublicKeyResponse
	err := doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.VpnIpAddresses, nil
}

func GetVpnNetworkDetails(installationToken string) (*VpnNetworkDetails, error) {
	requestBody := map[string]interface{}{
		"query": `
      query VpnNetworkDetails {
        vpnNetworkDetails {
          ip
          vpnEnabled
          geo {
            country
            city
            latitude
            longitude
          }
        }
      }
    `,
	}

	var response VpnNetworkDetailsResponse
	err := doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Details, nil
}

func GetVpnLocations(installationToken string) (*VpnLocations, error) {
	requestBody := map[string]interface{}{
		"query": `
      query VpnLocations {
        vpnLocations {
          countries {
            code
            name
            defaultCity
            cities {
              code
              unLocode
              name
              latitude
              longitude
              servers {
                hostname
                ipv4AddrIn
                ipv6AddrIn
                weight
                publicKey
                portRanges {
                  from
                  to
                }
                multihopPort
                ipv4Gateway
                ipv6Gateway
              }
            }
          }    
        }
      }
    `,
	}

	var response VpnLocationsResponse
	err := doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.VpnLocations, nil
}

func GetVpnClientDefaults(installationToken string) (*VpnClientDefaults, error) {
	requestBody := map[string]interface{}{
		"query": `
      query VpnClientDefaults {
        vpnClientDefaults {
          ipv4 {
            maskAll
            mtu
            maskNolocal
          }
          ipv6 {
            maskAll
            mtu
            maskNolocal
          }
          keyExpirationHours
        }
      }
    `,
	}

	var response VpnClientDefaultsResponse
	err := doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.VpnClientDefaults, nil
}

func doRequest(installationToken string, body *map[string]interface{}, responseBody any) error {
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", holocronStageUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	err = applyDefaultHeaders(installationToken, request)
  if err != nil {
    return err
  }

	if config.Debug {
		reqDump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("REQUEST:\n%s", string(reqDump))
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if config.Debug {
		respDump, err := httputil.DumpResponse(response, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("RESPONSE:\n%s", string(respDump))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

  var errorResponse ErrorResponse
  err = json.Unmarshal(data, &errorResponse)
  if err != nil {
    return err
  }
  if len(errorResponse.Errors) > 0 {
    return errors.New(errorResponse.Errors[0].Message)
  }

	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		return err
	}

	return nil
}

func applyDefaultHeaders(installationToken string, request *http.Request) error {
	request.Header.Set("Content-Type", "application/json")

	machineId, err := config.GetMachineId()
	if err != nil {
		return err
	}

	request.Header.Set("X-Device-Bearer", fmt.Sprintf("%s|%s", installationToken, machineId))

	return nil
}
