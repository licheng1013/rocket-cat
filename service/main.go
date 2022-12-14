package main

import (
	"core/core"
	"core/pkc"
)

func main() {
	ip := "192.168.101.10"
	service := core.NewService(ip,8999)
	defaultRpc := pkc.Grpc{}
	service.RpcLient(defaultRpc.RpcListen)
	service.Run(ip, 8848)
}
