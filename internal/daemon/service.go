package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	"github.com/malwarebytes/mbvpn-linux/pkg/connections"
	"github.com/malwarebytes/mbvpn-linux/pkg/remote"
	"github.com/malwarebytes/mbvpn-linux/pkg/servers"
	"github.com/malwarebytes/mbvpn-linux/pkg/session"
	"github.com/malwarebytes/mbvpn-linux/pkg/vpn"
)

type Service struct {
	Authorizer Authorizer
	Paths      Paths
	mu         sync.Mutex
}

func (s *Service) Handle(_ context.Context, caller rpc.Caller, request rpc.Request) (any, *rpc.Error) {
	if err := s.Authorizer.Authorize(caller); err != nil {
		return nil, &rpc.Error{Code: "not_authorized", Message: "caller is not authorized"}
	}
	switch request.Method {
	case "ping":
		return map[string]any{"protocol_version": rpc.Version}, nil
	case "login":
		var params struct {
			Code string `json:"code"`
		}
		if err := DecodeParams(request, &params); err != nil {
			return nil, err
		}
		return s.login(caller.UID, params.Code)
	case "logout":
		return s.logout(caller.UID)
	case "connect":
		var params struct {
			Selector string `json:"selector"`
		}
		if err := DecodeParams(request, &params); err != nil {
			return nil, err
		}
		return s.connect(caller.UID, params.Selector)
	case "disconnect":
		var params struct {
			Server   string `json:"server"`
			AllOwned bool   `json:"all_owned"`
		}
		if err := DecodeParams(request, &params); err != nil {
			return nil, err
		}
		return s.disconnect(caller.UID, params.Server, params.AllOwned)
	case "status":
		return s.status(caller.UID)
	case "locations":
		return s.locations(caller.UID)
	default:
		return nil, &rpc.Error{Code: "not_supported", Message: "unsupported operation"}
	}
}

func (s *Service) dependencies(uid uint32) (session.SessionManager, vpn.Vpn, error) {
	directory := config.NewDirectoryProvider(s.Paths.UserConfigDir(uid))
	configuration := config.NewYamlConfigProvider(directory)
	machineID := config.NewConfigFileMachineIdProvider(directory)
	holocron := remote.NewDefaultHolocron(machineID)
	vpnService, err := vpn.NewDefaultVpn(configuration, holocron, servers.NewDefaultServerStorage(directory), directory)
	if err != nil {
		return nil, nil, err
	}
	return session.NewDefaultSessionManager(configuration, holocron), vpnService, nil
}

func (s *Service) login(uid uint32, code string) (any, *rpc.Error) {
	sessionService, _, err := s.dependencies(uid)
	if err != nil {
		return nil, internalError(err)
	}
	if err := sessionService.LoginWithCode(code); err != nil {
		return nil, userError(err)
	}
	return map[string]any{"state": "active"}, nil
}

func (s *Service) logout(uid uint32) (any, *rpc.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sessionService, vpnService, err := s.dependencies(uid)
	if err != nil {
		return nil, internalError(err)
	}
	if err := vpnService.Disconnect(""); err != nil {
		return nil, userError(err)
	}
	if err := sessionService.Logout(); err != nil {
		return nil, userError(err)
	}
	return map[string]any{"state": "logged_out"}, nil
}

func (s *Service) connect(uid uint32, selector string) (any, *rpc.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if selector == "" {
		return nil, &rpc.Error{Code: "invalid_request", Message: "server selector is required"}
	}
	sessionService, vpnService, err := s.dependencies(uid)
	if err != nil {
		return nil, internalError(err)
	}
	if !sessionService.Active() {
		return nil, &rpc.Error{Code: "not_authenticated", Message: "there is no active session"}
	}
	if err := vpnService.Connect(selector); err != nil {
		return nil, userError(err)
	}
	return map[string]any{"state": "connected", "server": selector}, nil
}

func (s *Service) disconnect(uid uint32, server string, allOwned bool) (any, *rpc.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, vpnService, err := s.dependencies(uid)
	if err != nil {
		return nil, internalError(err)
	}
	if !allOwned && server == "" {
		return nil, &rpc.Error{Code: "invalid_request", Message: "connection target is required"}
	}
	if err := vpnService.Disconnect(server); err != nil {
		return nil, userError(err)
	}
	return map[string]any{"state": "disconnected"}, nil
}

func (s *Service) status(uid uint32) (any, *rpc.Error) {
	directory := config.NewDirectoryProvider(s.Paths.UserConfigDir(uid))
	store := connections.NewStore(directory)
	connections, err := store.List()
	if err != nil {
		return nil, internalError(err)
	}
	return map[string]any{"connections": connections}, nil
}

func (s *Service) locations(uid uint32) (any, *rpc.Error) {
	sessionService, vpnService, err := s.dependencies(uid)
	if err != nil {
		return nil, internalError(err)
	}
	if !sessionService.Active() {
		return nil, &rpc.Error{Code: "not_authenticated", Message: "there is no active session"}
	}
	if err := vpnService.Servers(true, true); err != nil {
		return nil, userError(err)
	}
	return map[string]any{"state": "updated"}, nil
}

func internalError(err error) *rpc.Error {
	return &rpc.Error{Code: "internal", Message: fmt.Sprintf("operation failed: %v", err)}
}
func userError(err error) *rpc.Error {
	return &rpc.Error{Code: "operation_failed", Message: err.Error()}
}

func DecodeParams(request rpc.Request, target any) *rpc.Error {
	if len(request.Params) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(request.Params))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return &rpc.Error{Code: "invalid_request", Message: "invalid request parameters"}
	}
	return nil
}
