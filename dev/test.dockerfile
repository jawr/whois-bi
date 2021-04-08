FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache git

COPY . .

RUN GOOS=linux \
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	go build -a -o toolbox main.go


FROM golang:alpine

RUN apk add netcat-openbsd bash

COPY --from=builder /build/toolbox /bin/toolbox

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

