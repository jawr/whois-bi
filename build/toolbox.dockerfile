FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux \
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	go build -a -o toolbox main.go

# actual service image
FROM busybox:latest

COPY --from=builder /build/toolbox /bin/toolbox

CMD ["bash"]
