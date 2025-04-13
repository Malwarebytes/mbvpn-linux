package console

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RunCmd executes a terminal command, optionally with sudo.
func RunCmd(sudo bool, command string, args ...string) (string, error) {
	var cmd *exec.Cmd

	var stdout strings.Builder
	var stderr strings.Builder
	
	if sudo {
		// Check if sudo is available
		_, err := exec.LookPath("sudo")
		if err != nil {
			log.Debugln("Sudo not available, running without sudo...")
			cmd = exec.Command(command, args...)
		} else {
			cmd = exec.Command("sudo", append([]string{command}, args...)...)
		}
	} else {
		cmd = exec.Command(command, args...)
	}	

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("command execution failed: %w, stderr: %v, stdout: %v", err, stderr.String(), stdout.String())
	}

	return stdout.String(), nil
}