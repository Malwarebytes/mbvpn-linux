package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
	"github.com/malwarebytes/mbvpn-linux/pkg/output"
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

func (c daemonClient) Servers(bool, bool) error { return c.call("locations", nil, nil) }
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

func (c daemonClient) call(method string, params any, result any) error {
	return c.client.Call(context.Background(), method, params, result)
}
