package core

import (
	"core/connect"
	"core/router"
	"log"
)

// Gateway 请使用 NewGateway 创建
type Gateway struct {
	// 设置网关连接方式
	socket connect.Socket
	// 单机模式可用
	router router.Router
	// 是否单机
	single bool
}

// SetRouter 设置路由器
func (g *Gateway) SetRouter(router router.Router) {
	g.router = router
}

// SetSingle 设置单机模式 true
func (g *Gateway) SetSingle(single bool) {
	g.single = single
}

func NewGateway() *Gateway {
	g := &Gateway{}
	g.SetSingle(true)
	g.SetRouter(router.Router{})
	return &Gateway{}
}

func (g *Gateway) Start(addr string, socket connect.Socket) {
	g.socket = socket
	g.socket.ListenBack(func(bytes []byte) []byte {
		log.Println("收到消息:" + string(bytes))
		return bytes
	})
	g.socket.ListenAddr(addr)
}
