# This how we want to name the binary output
GOBIN=go
BINARY=audify-rpc

BUILD=`git rev-parse --short HEAD`
SEMVERSION=1.0.0
VERSION="${SEMVERSION}@${BUILD}"
DEPSVERSION=`toml < Gopkg.lock | jq -c -r '.projects[] | "\(.name) \(.version)"' | tr ' ' ':' | tr '\n' ';' | sed 's/.$$//'`
UNAME_S := $(shell uname -s)
BUILD_OS=linux
ifeq ($(UNAME_S),Darwin)
    BUILD_OS=darwin
endif

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X github.com/theshadow/audify-rpc/cmd.BinaryVersion=${VERSION} -X github.com/theshadow/audify-rpc/cmd.BinaryDependencies=${DEPSVERSION}"

# Builds the project
build: deps
	CGO_ENABLED=0 GOOS=${BUILD_OS} go build ${LDFLAGS} -o ${BINARY}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

deps:
	dep ensure

test:
	go test ./...

.PHONY: clean install test
