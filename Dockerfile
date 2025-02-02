FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go build -o /mbvpn
CMD [ "/mbvpn" ]
