//go:build !linux

package daemon

import (
	"fmt"

	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
)

type Authorizer struct{ Group string }

func (a Authorizer) Authorize(rpc.Caller) error {
	return fmt.Errorf("mbvpnd is supported only on Linux")
}
