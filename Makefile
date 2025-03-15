include .env

BINARY_NAME=mbvpn

build-st:
	go build -o ${BINARY_NAME} -ldflags \
        "-X 'github.com/Malwarebytes/mbvpn/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL_ST}'" 

build-prod:
	go build -o ${BINARY_NAME} -ldflags \
        "-X 'github.com/Malwarebytes/mbvpn/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL_PROD}'"

run: install
	${BINARY_NAME}

clean:
	go clean
