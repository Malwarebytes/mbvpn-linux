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

	"github.com/Malwarebytes/mbvpn/pkg/config"
)

const (
	// baseUrl = "https://idms-holocron-stage.mwbsys.com/graphql"
	productCode      = "MBMA-C"
	productVersion   = "5.14.0"
)

var client = &http.Client{Timeout: time.Second * 10}

type Holocron interface {
	RegisterDevice() (string, error)
	ActivateDevice(installationToken string, key string, mbcode bool) (*DeviceModule, error)
	DeactivateDevice(installationToken string) (*DeviceModule, error)
	VpnRegisterPublicKey(installationToken string, key string) (*VpnIpAddresses, error)
	GetVpnNetworkDetails() (*VpnNetworkDetails, error)
	GetVpnLocations() (*VpnLocations, error)
	GetVpnClientDefaults() (*VpnClientDefaults, error)

	doRequest(installationToken string, body *map[string]interface{}, responseBody any) error
}

type DefaultHolocron struct {
	baseUrl     string
	mIdProvider config.MachineIdProvider
}

func NewDefaultHolocron(mIdProvider config.MachineIdProvider) Holocron {
	return &DefaultHolocron{
		baseUrl:     config.HolocronUrl,
		mIdProvider: mIdProvider,
	}
}

func (api *DefaultHolocron) RegisterDevice() (string, error) {
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

	var response HolocronResponse[RegisterDeviceResponseData]
	err := api.doRequest("", &requestBody, &response)
	if err != nil {
		return "", err
	}

	return response.Data.RegisterDevice.Device.InstallationToken, nil
}

func (api *DefaultHolocron) ActivateDevice(installationToken string, key string, mbcode bool) (*DeviceModule, error) {
	input := ActivateDeviceInput{
		Modules:        []ProductModule{ProductModulePrivacy},
		ActivationMode: ActivationModePassive,
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

	var response HolocronResponse[ActivateDeviceResponseData]
	err := api.doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.ActivateDevice.DeviceModules.Privacy, nil
}

func (api *DefaultHolocron) DeactivateDevice(installationToken string) (*DeviceModule, error) {
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

	var response HolocronResponse[DeactivateDeviceResponseData]
	err := api.doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.DeactivateDevice.DeviceModules.Privacy, nil
}

func (api *DefaultHolocron) VpnRegisterPublicKey(installationToken string, key string) (*VpnIpAddresses, error) {
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

	var response HolocronResponse[VpnRegisterPublicKeyResponseData]
	err := api.doRequest(installationToken, &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.VpnIpAddresses, nil
}

func (api *DefaultHolocron) GetVpnNetworkDetails() (*VpnNetworkDetails, error) {
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

	var response HolocronResponse[VpnNetworkDetailsResponseData]
	err := api.doRequest("", &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Details, nil
}

func (api *DefaultHolocron) GetVpnLocations() (*VpnLocations, error) {
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

	var response HolocronResponse[VpnLocationsResponseData]
	err := api.doRequest("", &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.VpnLocations, nil
}

func (api *DefaultHolocron) GetVpnClientDefaults() (*VpnClientDefaults, error) {
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

	var response HolocronResponse[VpnClientDefaultsResponseData]
	err := api.doRequest("", &requestBody, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.VpnClientDefaults, nil
}

func (api *DefaultHolocron) doRequest(installationToken string, body *map[string]interface{}, responseBody any) error {
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", api.baseUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	machineId, err := api.mIdProvider.Get()
	if err != nil {
		return err
	}

	request.Header.Set("X-Device-Bearer", fmt.Sprintf("%s|%s", installationToken, machineId))

	if config.Verbose() {
		reqDump, err := httputil.DumpRequest(request, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("REQUEST:\n%s\n", string(reqDump))
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if config.Verbose() {
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
