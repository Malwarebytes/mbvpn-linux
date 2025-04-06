FROM golang:1.24

RUN apt-get update && apt-get install -y \
    wireguard \
    iproute2 \
    iptables \
    net-tools \
    dbus \
    sudo \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV GOFLAGS="-buildvcs=false"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build-st

CMD ["sh", "-c", "dbus-uuidgen > /etc/machine-id && tail -f /dev/null"]
