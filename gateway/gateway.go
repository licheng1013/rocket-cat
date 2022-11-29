package gateway

import "io-game-go/core"

type Gateway struct {
	Nacos core.Nacos
}

func (n *Gateway) Register() {
	n.Nacos.SetServerConfig("192.168.101.10", 8848)
	n.Nacos.Register("192.168.101.10", 8001, core.GatewayName)
}

func (n *Gateway) GetServerInstance() {

}
