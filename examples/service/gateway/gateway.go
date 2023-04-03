package main

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
)

func main() {
	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: 12344, //这里是rpc端口
		ServiceName: common.GatewayName, RemoteName: common.ServiceName} //测试时 RemoteName 传递一样的
	nacos := registers.NewNacos()
	nacos.RegisterClient(clientInfo)
	nacos.Register(registers.RegisterInfo{Ip: "localhost", Port: 8848})

	gateway := core.DefaultGateway()
	gateway.SetSingle(false)
	gateway.SetSocket(&connect.WebSocket{})
	gateway.SetClient(&remote.GrpcClient{})
	gateway.SetRegisterClient(nacos)
	gateway.Start(connect.Addr)
}
