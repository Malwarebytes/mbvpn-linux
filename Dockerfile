FROM golang:1.23-alpine
WORKDIR /
COPY go.mod ./
RUN go mod download
COPY ../mbvpn /mbvpn
# RUN go build -o /mbvpn
CMD [ "ls" ]
