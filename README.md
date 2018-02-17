# Audify RPC

A gRPC wrapper for the audify.fm API. This was built to plug into a private microservices cluster. Feel free to use it however you want, you must follow the license.

## Installing 

`go get -u github.com/theshadow/audify-rpc`

## Building

### Locally

You'll need Golang 1.9.4 or greater installed. Then you'll want to run the make file `make build` you can also run the tests with `make tests`.

### Docker

Build the docker image with the `docker build -t audify-rpc:latest` command.

The conatainer will expose port **50051**.

## CLI

The service is also its own CLI tool. You can interact with your running instance using the `audify-rpc version` and `audify-rpc search` commands. 

## API

The gRPC interface can be found in the `service/` directory.