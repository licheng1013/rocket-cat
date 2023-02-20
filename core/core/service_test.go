package core

import (
	"github.com/io-game-go/common"
	"github.com/io-game-go/decoder"
	"github.com/io-game-go/registers"
	"github.com/io-game-go/remote"
	"log"
	"testing"
)

func TestService(t *testing.T) {
	log.Println("HelloWorld")
	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: 12345,
		ServiceName: common.ServicerName, RemoteName: common.ServicerName} // 测试时 RemoteName 传递一样的
	nacos := registers.NewNacos()
	nacos.RegisterClient(clientInfo)
	nacos.Register(registers.RegisterInfo{Ip: "localhost", Port: 8848})
	// nacos
	rpc := &remote.GrpcServer{}
	rpc.SetRegister(nacos)
	// rpc
	service := NewService()
	service.SetRpcServer(rpc)
	// 编码器
	service.SetDecoder(decoder.JsonDecoder{})

	// 关机钩子
	service.AddClose(nacos.Close)
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	service.Start()
}
