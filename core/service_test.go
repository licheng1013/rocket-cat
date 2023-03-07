package core

import (
	"github.com/licheng1013/io-game-go/common"
	"github.com/licheng1013/io-game-go/decoder"
	"github.com/licheng1013/io-game-go/messages"
	"github.com/licheng1013/io-game-go/protof"
	"github.com/licheng1013/io-game-go/registers"
	"github.com/licheng1013/io-game-go/remote"
	"github.com/licheng1013/io-game-go/router"
	"log"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	ports := []uint16{12000, 12001}
	//go ManyService(ports[0])
	ManyService(ports[1])
}

type MyProxy struct {
	proxy router.Proxy
}

func (m *MyProxy) InvokeFunc(ctx *router.Context) {
	log.Println("代理之前")
	m.proxy.InvokeFunc(ctx)
	log.Println("代理之后")
}

func (m *MyProxy) SetProxy(proxy router.Proxy) {
	m.proxy = proxy
}

// 此处测试需要配合注册中心一起测试
func ManyService(port uint16) {
	rpcClient := &remote.GrpcClient{}
	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: port,
		ServiceName: common.ServicerName, RemoteName: common.ServicerName} // 测试时 RemoteName 传递一样的
	nacos := registers.NewNacos()
	nacos.RegisterClient(clientInfo)
	nacos.Register(registers.RegisterInfo{Ip: "localhost", Port: 8848})
	// nacos
	rpc := &remote.GrpcServer{}
	// rpc
	service := NewService()
	service.SetRpcServer(rpc)
	service.SetRegister(nacos)
	// 编码器
	service.SetDecoder(decoder.JsonDecoder{})

	service.Router().AddProxy(&MyProxy{}) // 自定义注入器

	service.Router().AddAction(common.CmdKit.GetMerge(1, 2), func(ctx *router.Context) {
		ctx.Data = []byte("1")
	})
	service.Router().AddAction(common.CmdKit.GetMerge(1, 1), func(ctx *router.Context) {
		jsonMessage := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 2)}
		if ip, err := nacos.GetIp(); err == nil {
			bytes := rpcClient.InvokeRemoteRpc(ip.Addr(), protof.RpcBodyBuild(jsonMessage.GetBytesResult()))
			log.Println("调用逻辑服其他方法结果:", string(bytes))
		} else {
			log.Println("错误信息:" + err.Error())
		}
	})
	// 关机钩子
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	go service.Start()
	time.Sleep(5 * time.Second)
	jsonMessage := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1)}
	rpcClient.InvokeRemoteRpc(clientInfo.Addr(), protof.RpcBodyBuild(jsonMessage.GetBytesResult()))
	nacos.Close()
}
