include .env

BINARY_NAME=mbvpn

build-st:
	go build -o ${BINARY_NAME} -ldflags \
        "-X 'github.com/Malwarebytes/mbvpn/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL_ST}'" 

build-prod:
	go build -o ${BINARY_NAME} -ldflags \
        "-X 'github.com/Malwarebytes/mbvpn/pkg/config.HolocronUrl=${MBVPN_HOLOCRON_URL_PROD}'"

install:
	go install ./...

clean:
	go clean

docker-restart:
	docker-compose down
	docker-compose build
	docker-compose up -d
	docker-compose exec mbvpn bash
