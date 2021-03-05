FROM golang:alpine

WORKDIR /build

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy in the source code
COPY . .

RUN go install github.com/githubnemo/CompileDaemon

# arg is only available during build, so we need to 
# take a copy for runtime
ARG service
ENV build_service=$service

# we need real certs for email
RUN apk add ca-certificates

# use CompileDaemon to hot reload 
ENTRYPOINT CompileDaemon \
	-log-prefix=false \
	-build="go build -o /bin/service pkg/services/${build_service}/main.go" \
	-command="/bin/service"
