package cmd

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"strings"

	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
	"github.com/malwarebytes/mbvpn-linux/pkg/errors"
	"github.com/malwarebytes/mbvpn-linux/pkg/output"
	"github.com/malwarebytes/mbvpn-linux/pkg/remote"
)

type daemonClient struct{ client rpc.Client }

func newDaemonClient() daemonClient {
	socket := os.Getenv("MBVPN_SOCKET")
	if socket == "" {
		socket = "/run/mbvpn/mbvpnd.sock"
	}
	return daemonClient{client: rpc.Client{Socket: socket}}
}

func (c daemonClient) LoginWithCode(code string) error {
	return c.call("login", map[string]string{"code": code}, nil)
}

func (c daemonClient) Logout() error { return c.call("logout", nil, nil) }
func (c daemonClient) Active() bool  { return true }

func (c daemonClient) Servers(showCities bool, showServers bool) error {
	var result struct {
		Locations remote.VpnLocations `json:"locations"`
	}
	if err := c.call("locations", nil, &result); err != nil {
		return err
	}
	for _, country := range result.Locations.Countries {
		output.PrintMsg(fmt.Sprintf("%s, %s", country.Name, country.Code), output.MsgOutput)
		if !showCities {
			continue
		}
		for cityIndex, city := range country.Cities {
			cityPrefix := "└─"
			if cityIndex < len(country.Cities)-1 {
				cityPrefix = "├─"
			}
			output.PrintMsg(fmt.Sprintf("  %s %s, %s", cityPrefix, city.Name, city.Code), output.MsgOutput)
			if !showServers {
				continue
			}
			for serverIndex, server := range city.Servers {
				serverPrefix := "└─"
				if serverIndex < len(city.Servers)-1 {
					serverPrefix = "├─"
				}
				output.PrintMsg(fmt.Sprintf("     %s %s", serverPrefix, strings.Split(server.Hostname, ".")[0]), output.MsgOutput)
			}
		}
	}
	return nil
}
func (c daemonClient) Connect(selector string) error {
	return c.call("connect", map[string]string{"selector": selector}, nil)
}
func (c daemonClient) Disconnect(server string) error {
	params := map[string]any{"server": server, "all_owned": server == ""}
	return c.call("disconnect", params, nil)
}
func (c daemonClient) Status() error {
	var result struct {
		Connections []struct {
			ServerName    string `json:"server_name"`
			InterfaceName string `json:"interface_name"`
		} `json:"connections"`
	}
	if err := c.call("status", nil, &result); err != nil {
		return err
	}
	if len(result.Connections) == 0 {
		output.PrintMsg("No active connections.", output.MsgOutput)
		return nil
	}
	for _, connection := range result.Connections {
		output.PrintMsg(fmt.Sprintf("Connected to: %s (%s)", connection.ServerName, connection.InterfaceName), output.MsgSuccess)
	}
	return nil
}

func (daemonClient) Close() error { return nil }

func (c daemonClient) call(method string, params any, result any) error {
	err := c.client.Call(context.Background(), method, params, result)
	var rpcErr *rpc.Error
	if !stderrors.As(err, &rpcErr) {
		return err
	}
	return c.mapError(rpcErr)
}

func (daemonClient) mapError(rpcErr *rpc.Error) error {
	switch rpcErr.Code {
	case "not_authenticated":
		return errors.NewUserError(rpcErr.Message, errors.ErrUnauthorized)
	case "not_authorized":
		return errors.NewUserError(rpcErr.Message, errors.ErrPermissionDenied)
	case "invalid_request":
		return errors.NewUserError(rpcErr.Message, errors.ErrInvalidInput)
	case "operation_failed":
		return errors.NewUserError(rpcErr.Message, nil)
	default:
		return rpcErr
	}
}
