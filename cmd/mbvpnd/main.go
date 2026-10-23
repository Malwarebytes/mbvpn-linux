package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/malwarebytes/mbvpn-linux/internal/daemon"
	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
)

func main() {
	var socket string
	var group string
	flag.StringVar(&socket, "socket", "/run/mbvpn/mbvpnd.sock", "Unix socket path")
	flag.StringVar(&group, "group", "mbvpn", "authorized Unix group")
	flag.Parse()

	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "mbvpnd must run as root")
		os.Exit(1)
	}
	listener, err := rpc.Listen(socket, group)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	server := rpc.Server{Handler: &daemon.Service{Authorizer: daemon.Authorizer{Group: group}, Paths: daemon.Paths{StateDir: "/var/lib/mbvpn"}}}
	if err := server.Serve(listener, ctx); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}
}
