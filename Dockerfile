FROM docker.io/golang:1.24

RUN apt-get update && apt-get install -y \
    curl \
    sudo \
    wireguard \
    iproute2 \
    iptables \
    net-tools

RUN curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | sudo -E bash

RUN apt-get update && apt-get install -y task

RUN rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV GOFLAGS="-buildvcs=false"
ENV PATH="/app:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN task build-debug

CMD ["sh", "-c", "tail -f /dev/null"]
