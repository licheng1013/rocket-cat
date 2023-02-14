package core

import (
	"core/connect"
	"core/decoder"
	"core/router"
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

// Router 获取路由器
func (g *Gateway) Router() router.Router {
	return g.router
}

// SetRouter 设置自定义路由器->默认 router.DefaultRouter
func (g *Gateway) SetRouter(router router.Router) {
	g.router = router
}

// SetDecoder 设置编码器
func (g *Gateway) SetDecoder(decoder decoder.Decoder) {
	g.decoder = decoder
}

// SetSingle 设置单机模式 true
func (g *Gateway) SetSingle(single bool) {
	g.single = single
}

// NewGateway 默认以单机模式启动
func NewGateway() *Gateway {
	g := &Gateway{}
	g.SetSingle(true)
	g.router = &router.DefaultRouter{}
	return g
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
	if g.single {
		message := g.decoder.DecoderBytes(bytes)
		message.SetBody(g.router.InvokeFunc(message))
		if len(message.GetBody()) == 0 { // 没数据直接不返回
			return make([]byte, 0)
		}
		return message.GetBytesResult()
	}
	panic("描述: 目前暂时未实现rpc调用")
	return bytes
}
