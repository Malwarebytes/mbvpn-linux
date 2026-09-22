//go:build !linux

package rpc

import (
	"context"
	"fmt"
	"net"
)

type Caller struct {
	UID uint32
	GID uint32
	PID int32
}

type Handler interface {
	Handle(context.Context, Caller, Request) (any, *Error)
}

type Server struct{ Handler Handler }

func (s Server) Serve(net.Listener, context.Context) error {
	return fmt.Errorf("mbvpnd is supported only on Linux")
}

func Listen(string, string) (net.Listener, error) {
	return nil, fmt.Errorf("mbvpnd is supported only on Linux")
}
