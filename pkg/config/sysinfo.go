package config

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	cachedUserAgent string
	userAgentOnce   sync.Once
)

// GetKernelVersion returns the Linux kernel version using uname -r
func GetKernelVersion() string {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// GetDistroName returns the Linux distribution name
func GetDistroName() string {
	// Try to get distro name from lsb_release
	cmd := exec.Command("lsb_release", "-ds")
	output, err := cmd.Output()
	if err == nil {
		distro := strings.TrimSpace(string(output))
		// Remove quotes if present
		distro = strings.Trim(distro, "\"")
		return distro
	}

	// Fallback to uname -o
	cmd = exec.Command("uname", "-o")
	output, err = cmd.Output()
	if err != nil {
		return "Linux"
	}
	return strings.TrimSpace(string(output))
}

// GetArchitecture returns the CPU architecture
func GetArchitecture() string {
	return runtime.GOARCH
}

// GetUserAgent constructs the user agent string according to the template:
// Malwarebytes Privacy/$VERSION_NAME (github.com/malwarebytes/mbvpn-linux; Build: $VERSION_NAME; $LINUX_DISTRO_NAME $KERNEL_VERSION $ARCH)
func GetUserAgent(productVersion string) string {
	userAgentOnce.Do(func() {
		distro := GetDistroName()
		kernel := GetKernelVersion()
		arch := GetArchitecture()

		cachedUserAgent = fmt.Sprintf(
			"Malwarebytes Privacy/%s (github.com/malwarebytes/mbvpn-linux; Build: %s; %s %s %s)",
			productVersion,
			productVersion,
			distro,
			kernel,
			arch,
		)
	})

	return cachedUserAgent
}
