package core

import (
	"fmt"
	"github.com/io-game-go/common"
	"github.com/io-game-go/decoder"
	"github.com/io-game-go/registers"
	"github.com/io-game-go/remote"
	"github.com/io-game-go/router"
	"log"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	log.Println("HelloWorld")
	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: 12000,
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

	// 测试
	var count int64
	start := time.Now().UnixMilli()
	service.Router().AddFunc(10, func(ctx router.Context) []byte {
		ctx.Message.SetBody([]byte("Hi Ok"))
		end := time.Now().UnixMilli()
		count++
		if end-start > 1000 {
			fmt.Println("1s请求数量:", count)
			count = 0
			start = end
		}
		//log.Println(string(ctx.Message.GetBody()))
		return ctx.Message.GetBytesResult()
	})

	// 关机钩子
	service.AddClose(nacos.Close)
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	service.Start()
}
