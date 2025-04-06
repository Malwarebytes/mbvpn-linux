FROM golang:1.24

# Install wireguard and other dependencies
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

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Command to run when container starts
CMD ["sh", "-c", "dbus-uuidgen > /etc/machine-id && tail -f /dev/null"]
