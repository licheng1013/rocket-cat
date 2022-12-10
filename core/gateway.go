package core

import (
	"core/core"
)

type Gateway struct {
	Nacos  *core.Nacos
	Server *core.Server
}

func NewGateway() *Gateway {
	g := &Gateway{}
	g.Server = core.NewGameServer()
	g.Nacos = core.NewNacos()
	return g
}

// Run 注册中心地址和端口
func (n *Gateway) Run(ip string, port uint64) {
	n.Nacos.SetServerConfig(ip, port)
	n.Nacos.Register("192.168.101.10", n.Server.Port, core.GatewayName)
	n.Nacos.Init()      //初始化
	n.Nacos.Heartbeat() //心跳服务
	go func() {
		n.Server.Run()
	}()
}

func (n *Gateway) GetServerInstance() {

}
