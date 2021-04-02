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

# actual service image
FROM busybox:latest

COPY --from=builder /build/toolbox /bin/toolbox

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini
RUN chmod +x /tini

ENTRYPOINT ["/tini", "--"]
