#build stage
FROM golang:1.21.7-alpine AS builder
ARG SERVICE_NAME
WORKDIR /go/src
# ADD . .
# first copy modules that should be downloaded
COPY go.mod go.sum ./

ENV GO111MODULE=on

RUN GOPROXY=https://goproxy.io,direct go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o app cmd/$SERVICE_NAME/main.go

#final stage
FROM alpine:3.19.0
WORKDIR /srv
COPY --from=builder /go/src/ .
CMD ["./app"]

EXPOSE 5001