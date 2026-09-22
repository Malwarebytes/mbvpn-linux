//go:build linux

package daemon

import (
	"fmt"
	"os/user"
	"strconv"

	"github.com/malwarebytes/mbvpn-linux/internal/rpc"
)

type Authorizer struct {
	Group string
}

func (a Authorizer) Authorize(caller rpc.Caller) error {
	if caller.UID == 0 {
		return nil
	}
	group, err := user.LookupGroup(a.Group)
	if err != nil {
		return fmt.Errorf("lookup authorization group: %w", err)
	}
	account, err := user.LookupId(strconv.FormatUint(uint64(caller.UID), 10))
	if err != nil {
		return fmt.Errorf("lookup caller: %w", err)
	}
	groups, err := account.GroupIds()
	if err != nil {
		return fmt.Errorf("lookup caller groups: %w", err)
	}
	for _, id := range groups {
		if id == group.Gid {
			return nil
		}
	}
	return fmt.Errorf("caller is not a member of %s", a.Group)
}
