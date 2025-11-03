-include .env

BINARY_NAME=mbvpn
GO_PKG=github.com/malwarebytes/mbvpn-linux
GOFLAGS ?=

# Version information
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 0
VERSION_PATCH ?= 1
VERSION_BUILD ?= 0

# Common flags
COMMON_LDFLAGS=-X '${GO_PKG}/pkg/config.VersionMajor=${VERSION_MAJOR}' \
               -X '${GO_PKG}/pkg/config.VersionMinor=${VERSION_MINOR}' \
               -X '${GO_PKG}/pkg/config.VersionPatch=${VERSION_PATCH}' \
               -X '${GO_PKG}/pkg/config.VersionBuild=${VERSION_BUILD}'

# Environment flags
ENV_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildEnv=production' -X '${GO_PKG}/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL}'

# Debug flags
DEBUG_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildType=debug'
DEBUG_GCFLAGS=-N -l

# Release flags
RELEASE_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildType=release'
RELEASE_GCFLAGS=-trimpath

# Build targets
.PHONY: build build-debug build-release
build-debug:
	go build ${GOFLAGS} -o ${BINARY_NAME} -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ENV_LDFLAGS} ${DEBUG_LDFLAGS}"

build-release:
	go build ${GOFLAGS} -o ${BINARY_NAME} -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ENV_LDFLAGS} ${RELEASE_LDFLAGS}"

# Default target
build: build-release

# Install targets
.PHONY: install install-debug install-release
install-debug:
	go install -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ENV_LDFLAGS} ${DEBUG_LDFLAGS}" .

install-release:
	go install -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ENV_LDFLAGS} ${RELEASE_LDFLAGS}" .

# Default install target
install: install-release

# Test targets
.PHONY: test test-unit test-integration test-e2e clean
test:
	go test ./...

test-unit:
	go test -v ./pkg/... ./cmd/...

test-integration: build-release
	go test -tags=integration -v ./test/integration/...

test-e2e: build-release
	MBVPN_TEST_LICENSE_KEY=${MBVPN_TEST_LICENSE_KEY} go test -tags=e2e -v ./test/e2e/...

clean:
	go clean
