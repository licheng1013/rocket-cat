package main

import (
	"core/pkc"
)

func main() {
	ip := "192.168.101.10"
	prot := 8999
	p := &pkc.Grpc{}
	p.RpcListen(ip, uint64(prot))
	select {

	}
}
