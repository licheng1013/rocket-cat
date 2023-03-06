package remote

import (
	"log"
	"testing"
	"time"
)

func TestGrpcServer(t *testing.T) {
	const addr = "192.168.101.10:10002"
	server := GrpcServer{}
	server.CallbackResult(func(bytes []byte) []byte {
		log.Println("收到数据: ", string(bytes))
		return []byte("Hi")
	})
	go server.ListenAddr(addr)
	time.Sleep(time.Second)
	client := GrpcClient{}
	rpcResult := client.InvokeRemoteRpc(addr, []byte("HelloWorld"))
	log.Println(string(rpcResult))
}
