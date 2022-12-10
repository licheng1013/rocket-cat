package core

import (
	"core/common"
	"core/register"
)

// Service 构建服务
type Service struct {
	Nacos *register.Nacos
	App   *App
}

func NewService() *Service {
	g := &Service{}
	g.Nacos = register.NewNacos()
	g.App = NewGameServer(g.Nacos)
	return g
}

// Run 注册中心地址和端口
func (n *Service) Run(ip string, port uint64) {
	n.Nacos.SetServerConfig(ip, port)
	n.Nacos.Register("192.168.101.10", n.App.port, common.ServicerName)
	n.Nacos.Init()      //初始化
	n.Nacos.Heartbeat() //心跳服务
	n.App.SetStopFunc(func() {
		n.Nacos.Logout()
	})
	n.App.Run()
}
