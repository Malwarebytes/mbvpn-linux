//go:build linux

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

type Caller struct {
	UID uint32
	GID uint32
	PID int32
}

type Handler interface {
	Handle(context.Context, Caller, Request) (any, *Error)
}

type Server struct {
	Handler Handler
}

func (s Server) Serve(listener net.Listener, ctx context.Context) error {
	for {
		if unixListener, ok := listener.(*net.UnixListener); ok {
			_ = unixListener.SetDeadline(time.Now().Add(time.Second))
		}
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			return err
		}
		go s.serveConn(ctx, conn)
	}
}

func (s Server) serveConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	caller, err := peerCredentials(conn)
	if err != nil {
		return
	}
	request, err := ReadRequest(conn)
	if err != nil {
		_ = WriteResponse(conn, Response{ID: "invalid", Error: &Error{Code: "invalid_request", Message: "invalid request"}})
		return
	}
	result, rpcErr := s.Handler.Handle(ctx, caller, request)
	response := Response{ID: request.ID, Error: rpcErr}
	if rpcErr == nil && result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			response.Error = &Error{Code: "internal", Message: "internal error"}
		} else {
			response.Result = data
		}
	}
	_ = WriteResponse(conn, response)
}

func peerCredentials(conn net.Conn) (Caller, error) {
	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		return Caller{}, fmt.Errorf("expected Unix connection")
	}
	raw, err := unixConn.SyscallConn()
	if err != nil {
		return Caller{}, err
	}
	var credentials *syscall.Ucred
	var controlErr error
	err = raw.Control(func(fd uintptr) {
		credentials, controlErr = syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
	})
	if err != nil {
		return Caller{}, err
	}
	if controlErr != nil || credentials == nil {
		return Caller{}, fmt.Errorf("read peer credentials: %w", controlErr)
	}
	return Caller{UID: credentials.Uid, GID: credentials.Gid, PID: credentials.Pid}, nil
}

func Listen(socket string, groupID int) (net.Listener, error) {
	if err := os.MkdirAll("/run/mbvpn", 0o750); err != nil {
		return nil, err
	}
	if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, err
	}
	if err := os.Chown(socket, 0, groupID); err != nil {
		listener.Close()
		return nil, err
	}
	if err := os.Chmod(socket, 0o660); err != nil {
		listener.Close()
		return nil, err
	}
	return listener, nil
}
