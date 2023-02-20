package core

import (
	"github.com/io-game-go/common"
	"github.com/io-game-go/registers"
	"github.com/io-game-go/remote"
	"log"
	"testing"
)

func TestService(t *testing.T) {
	log.Println("HelloWorld")

	nacos := registers.NewNacos()
	nacos.Register(registers.RegisterInfo{Ip: "localhost", Port: 8848,
		ServiceName: common.ServicerName, RemoteName: common.GatewayName})

	// nacos
	rpc := &remote.GrpcServer{}
	rpc.SetRegister(nacos)

	// rpc
	service := NewService()
	service.SetRpcServer(rpc)

	// 关机
	service.AddClose(rpc.Close)
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	service.Start()
}
