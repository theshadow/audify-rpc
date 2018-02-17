// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>

//go:generate protoc -I service/ service/server.proto --go_out=plugins=grpc:service

package main

import "github.com/theshadow/audify-rpc/cmd"

func main() {
	cmd.Execute()
}
