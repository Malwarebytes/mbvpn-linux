package daemon

import (
	"fmt"
	"path/filepath"
	"strconv"
)

type Paths struct {
	StateDir string
}

func (p Paths) UserConfigDir(uid uint32) string {
	return filepath.Join(p.StateDir, "users", strconv.FormatUint(uint64(uid), 10))
}

func (p Paths) Validate() error {
	if !filepath.IsAbs(p.StateDir) {
		return fmt.Errorf("state directory must be absolute")
	}
	return nil
}
