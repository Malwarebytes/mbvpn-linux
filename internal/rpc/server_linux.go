//go:build linux

package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	maxActiveConnections = 32
	connectionTimeout    = 30 * time.Second
)

var ErrSocketInUse = errors.New("socket is already in use")

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
	slots := make(chan struct{}, maxActiveConnections)
	for {
		if ctx.Err() != nil {
			return nil
		}
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
		select {
		case slots <- struct{}{}:
		case <-ctx.Done():
			conn.Close()
			return nil
		}
		go func() {
			defer func() { <-slots }()
			s.serveConn(ctx, conn)
		}()
	}
}

func (s Server) serveConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()
	stopOnCancel := context.AfterFunc(ctx, func() {
		_ = conn.Close()
	})
	defer stopOnCancel()
	deadline, _ := ctx.Deadline()
	_ = conn.SetDeadline(deadline)
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
	if err = WriteResponse(conn, response); errors.Is(err, ErrFrameTooLarge) {
		_ = WriteResponse(conn, Response{
			ID:    request.ID,
			Error: &Error{Code: "response_too_large", Message: "response exceeds maximum size"},
		})
	}
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
	var listener net.Listener
	err := withSocketLock("/run/mbvpn/daemon.lock", func() error {
		if err := prepareSocket(socket); err != nil {
			return err
		}
		if err := os.Chown("/run/mbvpn", 0, groupID); err != nil {
			return err
		}
		if err := os.Chmod("/run/mbvpn", 0o750); err != nil {
			return err
		}
		ln, err := net.Listen("unix", socket)
		if err != nil {
			return err
		}
		if err := os.Chown(socket, 0, groupID); err != nil {
			ln.Close()
			_ = os.Remove(socket)
			return err
		}
		if err := os.Chmod(socket, 0o660); err != nil {
			ln.Close()
			_ = os.Remove(socket)
			return err
		}
		listener = ln
		return nil
	})
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func prepareSocket(socket string) error {
	conn, err := net.DialTimeout("unix", socket, time.Second)
	if err == nil {
		conn.Close()
		return ErrSocketInUse
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		info, statErr := os.Stat(socket)
		if errors.Is(statErr, os.ErrNotExist) {
			return nil
		}
		if statErr != nil {
			return statErr
		}
		if info.Mode()&os.ModeSocket != 0 {
			return os.Remove(socket)
		}
	}
	return err
}

func withSocketLock(lockPath string, fn func() error) error {
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o640)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	return fn()
}
