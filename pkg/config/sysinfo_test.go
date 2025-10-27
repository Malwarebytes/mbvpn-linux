package config

import (
	"strings"
	"testing"
)

func TestGetKernelVersion(t *testing.T) {
	version := GetKernelVersion()
	if version == "" {
		t.Error("Expected non-empty kernel version")
	}
	if version == "unknown" {
		t.Log("Warning: Could not determine kernel version")
	}
	t.Logf("Kernel version: %s", version)
}

func TestGetDistroName(t *testing.T) {
	distro := GetDistroName()
	if distro == "" {
		t.Error("Expected non-empty distro name")
	}
	t.Logf("Distro name: %s", distro)
}

func TestGetArchitecture(t *testing.T) {
	arch := GetArchitecture()
	if arch == "" {
		t.Error("Expected non-empty architecture")
	}
	// Common architectures
	validArchs := []string{"amd64", "386", "arm64", "arm", "ppc64le", "s390x"}
	isValid := false
	for _, valid := range validArchs {
		if arch == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		t.Logf("Warning: Unusual architecture: %s", arch)
	}
	t.Logf("Architecture: %s", arch)
}

func TestGetUserAgent(t *testing.T) {
	productVersion := "5.14.0"
	userAgent := GetUserAgent(productVersion)

	if userAgent == "" {
		t.Fatal("Expected non-empty user agent")
	}

	// Verify format: Malwarebytes Privacy/VERSION (github.com/Malwarebytes/mbvpn-linux; Build: VERSION; DISTRO KERNEL ARCH)
	if !strings.HasPrefix(userAgent, "Malwarebytes Privacy/"+productVersion) {
		t.Errorf("Expected user agent to start with 'Malwarebytes Privacy/%s', got: %s", productVersion, userAgent)
	}

	if !strings.Contains(userAgent, "github.com/Malwarebytes/mbvpn-linux") {
		t.Error("Expected user agent to contain 'github.com/Malwarebytes/mbvpn-linux'")
	}

	if !strings.Contains(userAgent, "Build: "+productVersion) {
		t.Errorf("Expected user agent to contain 'Build: %s'", productVersion)
	}

	t.Logf("User agent: %s", userAgent)

	// Test caching - should return same value
	userAgent2 := GetUserAgent(productVersion)
	if userAgent != userAgent2 {
		t.Error("Expected cached user agent to match first call")
	}
}

func TestGetUserAgentFormat(t *testing.T) {
	productVersion := "5.14.0"
	userAgent := GetUserAgent(productVersion)

	// Verify structure
	parts := strings.Split(userAgent, " (")
	if len(parts) != 2 {
		t.Fatalf("Expected user agent to have format 'PREFIX (DETAILS)', got: %s", userAgent)
	}

	prefix := parts[0]
	if prefix != "Malwarebytes Privacy/"+productVersion {
		t.Errorf("Expected prefix 'Malwarebytes Privacy/%s', got: %s", productVersion, prefix)
	}

	details := strings.TrimSuffix(parts[1], ")")
	detailParts := strings.Split(details, "; ")
	if len(detailParts) != 3 {
		t.Errorf("Expected 3 detail parts separated by '; ', got %d parts: %v", len(detailParts), detailParts)
	}

	// Check first detail part
	if detailParts[0] != "github.com/Malwarebytes/mbvpn-linux" {
		t.Errorf("Expected first detail 'github.com/Malwarebytes/mbvpn-linux', got: %s", detailParts[0])
	}

	// Check second detail part
	if !strings.HasPrefix(detailParts[1], "Build: ") {
		t.Errorf("Expected second detail to start with 'Build: ', got: %s", detailParts[1])
	}

	// Check third detail part has distro, kernel, and arch
	systemInfo := detailParts[2]
	systemParts := strings.Fields(systemInfo)
	if len(systemParts) < 2 {
		t.Errorf("Expected system info to have at least distro and kernel, got: %s", systemInfo)
	}
}
