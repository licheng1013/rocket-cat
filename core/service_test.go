package core

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
	"github.com/licheng1013/rocket-cat/router"
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
	service.SetRpcClient(&remote.GrpcClient{})

	service.Router().AddProxy(&MyProxy{}) // 自定义注入器

	service.Router().AddAction(common.CmdKit.GetMerge(1, 2), func(ctx *router.Context) {
		ctx.Data = []byte("1")
	})
	service.Router().AddAction(common.CmdKit.GetMerge(1, 1), func(ctx *router.Context) {
		jsonMessage := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 2)}
		message, err := service.SendMessage(jsonMessage.GetBytesResult())
		if err != nil {
			log.Println("错误:", err.Error())
			return
		}
		for _, item := range message {
			log.Println(string(item))
		}
	})
	// 关机钩子
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	go service.Start()
	time.Sleep(5 * time.Second)
	jsonMessage := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1)}
	rpcClient := &remote.GrpcClient{}
	rpcClient.InvokeRemoteRpc(clientInfo.Addr(), protof.RpcBodyBuild(jsonMessage.GetBytesResult()))
	nacos.Close()
}
