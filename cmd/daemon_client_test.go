package cmd

import (
	"errors"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
	mbvpnerors "github.com/malwarebytes/mbvpn-linux/pkg/errors"
)

func TestDaemonErrorMapsToUserError(t *testing.T) {
	client := daemonClient{}
	err := client.mapError(&rpc.Error{Code: "not_authenticated", Message: "there is no active session"})
	if !mbvpnerors.IsUserError(err) || !mbvpnerors.IsUnauthorized(err) {
		t.Fatalf("expected an unauthorized user error, got %v", err)
	}
}

func TestDaemonErrorLeavesInternalErrorsTyped(t *testing.T) {
	client := daemonClient{}
	err := client.mapError(&rpc.Error{Code: "internal", Message: "operation failed"})
	var rpcErr *rpc.Error
	if !errors.As(err, &rpcErr) || rpcErr.Code != "internal" {
		t.Fatalf("expected typed RPC error, got %v", err)
	}
}
