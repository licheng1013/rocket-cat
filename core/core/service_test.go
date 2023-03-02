package core

import (
	"fmt"
	"github.com/io-game-go/common"
	"github.com/io-game-go/decoder"
	"github.com/io-game-go/registers"
	"github.com/io-game-go/remote"
	"github.com/io-game-go/router"
	"log"
	"sync"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	ports := []uint16{12000, 12001}
	go ManyService(ports[0])
	ManyService(ports[1])
}

func ManyService(port uint16) {
	var lock sync.Mutex
	log.Println("HelloWorld")
	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: port,
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

	service.Router().AddFunc(common.CmdKit.GetMerge(1, 2), func(ctx *router.Context) {
		ctx.Message.SetBody([]byte("Hello")).GetBytesResult()
	})

	service.Router().AddFunc(common.CmdKit.GetMerge(1, 1), func(ctx *router.Context) {
		ctx.Message.SetBody([]byte("Hi Ok"))
		end := time.Now().UnixMilli()
		lock.Lock()
		count++
		if end-start > 1000 {

			server := ctx.RpcServer.(*remote.GrpcServer)
			server.CountRoom()

			fmt.Println(port, "1s请求数量:", count)
			count = 0
			start = end
		}
		lock.Unlock()
		//log.Println(string(ctx.Message.GetBody()))
	})

	// 关机钩子
	service.AddClose(nacos.Close)
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	service.Start()
}
