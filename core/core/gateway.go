package core

import (
	"core/common"
	"core/register"
	"fmt"
	"log"
)

type Gateway struct {
	Nacos *register.Nacos
	App   *App
}

func NewGateway(ip string,port uint64) *Gateway {
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	g := &Gateway{}
	g.Nacos = register.NewNacos()
	g.App = NewGameServer(g.Nacos)
	g.App.SetIp(ip)
	g.App.SetProt(port)
	return g
}

// Run 注册中心地址和端口
func (n *Gateway) Run(nacosIp string, port uint64) {
	log.Println("nacos注册地址: http://"+ n.App.ip +":"+fmt.Sprint(port))
	log.Println("gateway注册地址: http://"+ nacosIp +":"+fmt.Sprint(n.App.port))
	n.Nacos.SetServerConfig(n.App.ip, port)
	n.Nacos.Register(nacosIp, n.App.port, common.GatewayName)
	n.Nacos.Init()      //初始化
	n.Nacos.Heartbeat() //心跳服务
	n.App.SetStopFunc(func() {
		n.Nacos.Logout()
	})
	n.App.Run()

}
