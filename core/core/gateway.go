package core

import (
	"core/connect"
	"core/decoder"
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
	// 编码器
	decoder decoder.Decoder
}

func (g *Gateway) SetDecoder(decoder decoder.Decoder) {
	g.decoder = decoder
}

func (g *Gateway) Router() router.Router {
	return g.router
}

func (g *Gateway) Decoder() decoder.Decoder {
	return g.decoder
}

// SetSingle 设置单机模式 true
func (g *Gateway) SetSingle(single bool) {
	g.single = single
}

func NewGateway() *Gateway {
	g := &Gateway{}
	g.SetSingle(true)
	g.router = router.Router{}
	return &Gateway{}
}

func (g *Gateway) Start(addr string, socket connect.Socket) {
	if g.decoder == nil {
		panic("没有设置编码器: decoder")
	}
	g.socket = socket
	g.socket.ListenBack(g.ListenBack)
	g.socket.ListenAddr(addr)
}

func (g *Gateway) ListenBack(bytes []byte) []byte {
	message := g.decoder.DecoderBytes(bytes)
	invoke := g.router.RouterMap[message.GetMerge()]
	if invoke == nil {
		log.Println("未注册路由方法:", message.GetMerge())
		return make([]byte, 0)
	}
	log.Println("收到消息:", string(message.GetBody()))
	return bytes
}
