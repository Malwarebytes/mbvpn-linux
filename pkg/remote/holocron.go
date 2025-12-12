package remote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	log "github.com/sirupsen/logrus"
)

const (
	productCode    = "MBMA-C"
	productVersion = "5.14.0"
)

var ErrorServerResponse = errors.New("error response from server")

var client = &http.Client{Timeout: time.Second * 10}

type Holocron interface {
	RegisterDevice() (string, error)
	ActivateDevice(installationToken string, code string) (*DeviceModule, error)
	CheckDevice(installationToken string) (*DeviceModule, error)
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
		baseUrl:     config.GetHolocronUrl(),
		mIdProvider: mIdProvider,
	}
}

func (api *DefaultHolocron) RegisterDevice() (string, error) {
	log.Debugln("Begin: `RegisterDevice` request")
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
		log.Debugln("Error: `RegisterDevice` request")
		return "", fmt.Errorf("failed to register device: %w", err)
	}

	log.Debugln("Success: `RegisterDevice` request")
	return response.Data.RegisterDevice.Device.InstallationToken, nil
}

func (api *DefaultHolocron) ActivateDevice(installationToken string, code string) (*DeviceModule, error) {
	log.Debugln("Begin: `ActivateDevice` request")
	input := ActivateDeviceInput{
		Modules:          []ProductModule{ProductModulePrivacy},
		ActivationMode:   ActivationModePassive,
		ActivationMethod: ActivationMethodOneTimeToken,
		OneTimeToken:     code,
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
		log.Debugln("Error: `ActivateDevice` request")
		return nil, fmt.Errorf("failed to activate device: %w", err)
	}

	log.Debugln("Success: `ActivateDevice` request")
	return &response.Data.ActivateDevice.DeviceModules.Privacy, nil
}

func (api *DefaultHolocron) CheckDevice(installationToken string) (*DeviceModule, error) {
	log.Debugln("Begin: `CheckDevice` request")
	input := CheckDeviceInput{
		Modules: []ProductModule{ProductModulePrivacy},
	}

	requestBody := map[string]interface{}{
		"query": `
      mutation CheckDevice($input: CheckDeviceInput!) {
        checkDevice(input: $input) {
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

	var response HolocronResponse[CheckDeviceResponseData]
	err := api.doRequest(installationToken, &requestBody, &response)
	if err != nil {
		log.Debugln("Error: `CheckDevice` request")
		return nil, fmt.Errorf("failed to check device: %w", err)
	}

	log.Debugln("Success: `CheckDevice` request")
	return &response.Data.CheckDevice.DeviceModules.Privacy, nil
}

func (api *DefaultHolocron) DeactivateDevice(installationToken string) (*DeviceModule, error) {
	log.Debugln("Begin: `DeactivateDevice` request")
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
		log.Debugln("Error: `DeactivateDevice` request")
		return nil, fmt.Errorf("failed to deactivate device: %w", err)
	}

	log.Debugln("Success: `DeactivateDevice` request")
	return &response.Data.DeactivateDevice.DeviceModules.Privacy, nil
}

func (api *DefaultHolocron) VpnRegisterPublicKey(installationToken string, key string) (*VpnIpAddresses, error) {
	log.Debugln("Begin: `VpnRegisterPublicKey` request")
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
		log.Debugln("Error: `VpnRegisterPublicKey` request")
		return nil, fmt.Errorf("failed to register VPN public key: %w", err)
	}

	log.Debugln("Success: `VpnRegisterPublicKey` request")
	return &response.Data.VpnIpAddresses, nil
}

func (api *DefaultHolocron) GetVpnNetworkDetails() (*VpnNetworkDetails, error) {
	log.Debugln("Begin: `GetVpnNetworkDetails` request")
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
		log.Debugln("Error: `GetVpnNetworkDetails` request")
		return nil, fmt.Errorf("failed to get VPN network details: %w", err)
	}

	log.Debugln("Success: `GetVpnNetworkDetails` request")
	return &response.Data.Details, nil
}

func (api *DefaultHolocron) GetVpnLocations() (*VpnLocations, error) {
	log.Debugln("Begin: `GetVpnLocations` request")
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
		log.Debugln("Error: `GetVpnLocations` request")
		return nil, fmt.Errorf("failed to get VPN locations: %w", err)
	}

	log.Debugln("Success: `GetVpnLocations` request")
	return &response.Data.VpnLocations, nil
}

func (api *DefaultHolocron) GetVpnClientDefaults() (*VpnClientDefaults, error) {
	log.Debugln("Begin: `GetVpnClientDefaults` request")
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
		log.Debugln("Error: `GetVpnClientDefaults` request")
		return nil, fmt.Errorf("failed to get VPN client defaults: %w", err)
	}

	log.Debugln("Success: `GetVpnClientDefaults` request")
	return &response.Data.VpnClientDefaults, nil
}

func (api *DefaultHolocron) doRequest(installationToken string, body *map[string]interface{}, responseBody any) error {
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	request, err := http.NewRequest("POST", api.baseUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	machineId, err := api.mIdProvider.Get()
	if err != nil {
		return fmt.Errorf("failed to get machine id: %w", err)
	}

	request.Header.Set("X-Device-Bearer", fmt.Sprintf("%s|%s", installationToken, machineId))
	request.Header.Set("User-Agent", config.GetUserAgent(productVersion))

	reqDump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Errorln(fmt.Errorf("failed to generate request dump: %w", err))
	} else {
		log.Traceln("Request dump:", string(reqDump))
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	respDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Errorln(fmt.Errorf("failed to generate response dump: %w", err))
	} else {
		log.Traceln("Response dump:", string(respDump))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var errorResponse ErrorResponse
	err = json.Unmarshal(data, &errorResponse)
	if err != nil {
		return fmt.Errorf("failed to parse error response: %w", err)
	}
	if len(errorResponse.Errors) > 0 {
		return fmt.Errorf("backend error: %w", errors.Join(ErrorServerResponse, errors.New(errorResponse.Errors[0].Message)))
	}

	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	return nil
}
