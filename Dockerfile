FROM golang:1.9.4
LABEL description="audify-rpc provides a gRPC wrapper for the audify.fm API"
WORKDIR /go/src/github.com/theshadow/audify-rpc/
COPY ./ /go/src/github.com/theshadow/audify-rpc/
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN curl -sL https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-386 -o /usr/local/bin/dep
RUN chmod +x /usr/local/bin/dep
RUN apt-get update && apt-get install -y --force-yes nodejs jq
RUN npm install -g toml-cli
RUN make deps && make test && make build

FROM alpine:latest
EXPOSE 50051
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/theshadow/audify-rpc/audify-rpc .
CMD ["./audify-rpc", "start"]
