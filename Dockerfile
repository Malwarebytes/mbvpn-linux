FROM docker.io/golang:1.24

RUN apt-get update && apt-get install -y \
    wireguard \
    iproute2 \
    iptables \
    net-tools \
    sudo \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV GOFLAGS="-buildvcs=false"
ENV PATH="/app:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build-debug

CMD ["sh", "-c", "tail -f /dev/null"]
