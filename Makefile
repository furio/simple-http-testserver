BINARY=http-test-server
BINARY_PATH=dist
VERSION=$(shell git log --pretty=format:'%h' -n 1)
BUILD=$(shell date +%FT%T%z)
HUMAN_VERSION=$(shell cat VERSION)-${VERSION}

# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -tags netgo -o ${BINARY_PATH}/${BINARY}-${GOOS} simple-http-testserver.go
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X main.Version=${HUMAN_VERSION} -X main.Build=${BUILD}"

init:
	go mod vendor
	
build:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_PATH}/${BINARY}-darwin simple-http-testserver.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_PATH}/${BINARY}-linux simple-http-testserver.go

clean:
	rm -rf ${BINARY_PATH}/

do: clean init build
