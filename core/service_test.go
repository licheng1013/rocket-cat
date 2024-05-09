package core

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
	"github.com/licheng1013/rocket-cat/router"
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
	common.CatLog.Println("代理之前")
	m.proxy.InvokeFunc(ctx)
	common.CatLog.Println("代理之后")
}

func (m *MyProxy) SetProxy(proxy router.Proxy) {
	m.proxy = proxy
}

// 此处测试需要配合注册中心一起测试
func ManyService(port uint16) {
	clientInfo := registers.ClientInfo{Ip: "192.168.101.100", Port: port,
		ServiceName: common.ServiceName, RemoteName: common.ServiceName} // 测试时 RemoteName 传递一样的
	//register := registers.NewNacos()
	//register.RegisterClient(clientInfo)
	//register.RegisterServer(registers.ServerInfo{Ip: "localhost", Port: 8848})

	register := registers.NewEtcd()
	register.RegisterClient(clientInfo)
	register.RegisterServer(registers.ServerInfo{Ip: "localhost", Port: 2379})

	// register
	rpc := &remote.GrpcServer{}
	// rpc
	service := NewService()
	service.SetRpcServer(rpc)
	service.SetRegister(register)
	// 编码器
	service.SetDecoder(decoder.JsonDecoder{})
	service.SetRpcClient(&remote.GrpcClient{})

	service.Router().AddProxy(&MyProxy{}) // 自定义注入器

	service.Router().Action(1, 2, func(ctx *router.Context) {
		ctx.Result([]byte("1"))
	})
	service.Router().Action(1, 1, func(ctx *router.Context) {
		jsonMessage := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 2)}
		message, err := service.SendServiceMessage(jsonMessage.GetBytesResult())
		if err != nil {
			common.CatLog.Println("错误:", err.Error())
			return
		}
		for _, item := range message {
			common.CatLog.Println(string(item))
		}
	})
	// 关机钩子
	service.AddClose(func() {
		common.CatLog.Println("在关机中了")
	})
	go service.Start()
	time.Sleep(5 * time.Second)
	jsonMessage := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1)}
	rpcClient := &remote.GrpcClient{}
	rpcClient.InvokeRemoteRpc(clientInfo.Addr(), protof.RpcBodyBuild(jsonMessage.GetBytesResult()))
	register.Close()
}
