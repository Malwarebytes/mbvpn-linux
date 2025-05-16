include .env

BINARY_NAME=mbvpn
GO_PKG=github.com/Malwarebytes/mbvpn

# Common flags
COMMON_LDFLAGS=

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
	go build -o ${BINARY_NAME} -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ST_ENV_LDFLAGS} ${DEBUG_LDFLAGS}"

build-st-release:
	go build -o ${BINARY_NAME} -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${ST_ENV_LDFLAGS} ${RELEASE_LDFLAGS}"

build-prod-debug:
	go build -o ${BINARY_NAME} -gcflags="${DEBUG_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${PROD_ENV_LDFLAGS} ${DEBUG_LDFLAGS}"

build-prod-release:
	go build -o ${BINARY_NAME} -gcflags="${RELEASE_GCFLAGS}" -ldflags "${COMMON_LDFLAGS} ${PROD_ENV_LDFLAGS} ${RELEASE_LDFLAGS}"

# Default targets
build-st: build-st-release

build-prod: build-prod-release

install:
	go install ./...

clean:
	go clean

docker-restart:
	docker-compose down
	docker-compose build
	docker-compose up -d
	docker-compose exec mbvpn bash