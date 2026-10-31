//go:build !linux

package wireguard

import "fmt"

func NewManager() (Manager, error) {
	return nil, fmt.Errorf("WireGuard management is supported only on Linux")
}
