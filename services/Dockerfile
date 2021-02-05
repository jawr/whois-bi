FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG service
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -o service services/${service}/main.go

# actual service image
FROM busybox:latest

WORKDIR /srv

COPY --from=builder /build/service .

CMD ["./service"]
