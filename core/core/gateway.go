package core

import "log"

type Gateway struct {
	Nacos *Nacos
	App   *App
}

func NewGateway() *Gateway {
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	g := &Gateway{}
	g.App = NewGameServer()
	g.Nacos = NewNacos()
	return g
}

// Run 注册中心地址和端口
func (n *Gateway) Run(ip string, port uint64) {
	n.Nacos.SetServerConfig(ip, port)
	n.Nacos.Register("192.168.101.10", n.App.Port, GatewayName)
	n.Nacos.Init()      //初始化
	n.Nacos.Heartbeat() //心跳服务
	n.App.Run()
}
