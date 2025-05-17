-include .env

BINARY_NAME=mbvpn
GO_PKG=github.com/Malwarebytes/mbvpn
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
ST_ENV_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildEnv=staging' -X '${GO_PKG}/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL_ST}'
PROD_ENV_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildEnv=production' -X '${GO_PKG}/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL_PROD}'

# Debug flags
DEBUG_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildType=debug'
DEBUG_GCFLAGS=-N -l

# Release flags
RELEASE_LDFLAGS=-X '${GO_PKG}/pkg/config.BuildType=release'
RELEASE_GCFLAGS=-trimpath

# Build matrix targets
build-st-debug:
	go build ${GOFLAGS} -o ${BINARY_NAME} -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ST_ENV_LDFLAGS} ${DEBUG_LDFLAGS}"

build-st-release:
	go build ${GOFLAGS} -o ${BINARY_NAME} -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ST_ENV_LDFLAGS} ${RELEASE_LDFLAGS}"

build-prod-debug:
	go build ${GOFLAGS} -o ${BINARY_NAME} -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${PROD_ENV_LDFLAGS} ${DEBUG_LDFLAGS}"

build-prod-release:
	go build ${GOFLAGS} -o ${BINARY_NAME} -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${PROD_ENV_LDFLAGS} ${RELEASE_LDFLAGS}"

# Default targets
build-st: build-st-release

build-prod: build-prod-release

# Install targets
install-st-debug: 
	go install -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ST_ENV_LDFLAGS} ${DEBUG_LDFLAGS}" .

install-st-release: 
	go install -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ST_ENV_LDFLAGS} ${RELEASE_LDFLAGS}" .

install-prod-debug: 
	go install -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${PROD_ENV_LDFLAGS} ${DEBUG_LDFLAGS}" .

install-prod-release: 
	go install -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${PROD_ENV_LDFLAGS} ${RELEASE_LDFLAGS}" .

# Default install targets
install-st: install-st-release

install-prod: install-prod-release

install: install-prod

# Test targets
test:
	go test ./...

test-unit:
	go test -v ./pkg/... ./cmd/...

test-integration: build-st-release
	go test -tags=integration -v ./test/integration/...

test-e2e: build-st-release
	go test -tags=e2e -v ./test/e2e/...

# Installation script tests
test-install:
	./test/install/run_tests.sh

test-install-ubuntu:
	./test/install/run_tests.sh -d ubuntu

test-install-fedora:
	./test/install/run_tests.sh -d fedora

test-install-centos:
	./test/install/run_tests.sh -d centos

test-install-arch:
	./test/install/run_tests.sh -d arch

test-install-opensuse:
	./test/install/run_tests.sh -d opensuse

clean:
	go clean

docker-restart:
	docker-compose down
	docker-compose build
	docker-compose up -d
	docker-compose exec mbvpn bash
