package core

import (
	"github.com/licheng1013/io-game-go/common"
	"github.com/licheng1013/io-game-go/connect"
	"github.com/licheng1013/io-game-go/decoder"
	"github.com/licheng1013/io-game-go/registers"
	"github.com/licheng1013/io-game-go/remote"
	router "github.com/licheng1013/io-game-go/router"
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
	// Remote Invoke Client
	client remote.RpcClient
	// Register Client
	registerClient registers.Register
}

func (g *Gateway) SetClient(client remote.RpcClient) {
	g.client = client
}

func (g *Gateway) SetRegisterClient(registerClient registers.Register) {
	g.registerClient = registerClient
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
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	common.AssertNil(socket, "没有设置链接协议")
	if g.single {
		common.AssertNil(g.decoder, "没有设置编码器: decoder")
	} else {
		common.AssertNil(g.client, "没有设置远程客户端!")
		common.AssertNil(g.registerClient, "没有设置注册客户端!")
	}
	g.socket = socket
	g.socket.ListenBack(g.ListenBack)
	log.Println("监听Socket:" + addr)
	go g.socket.ListenAddr(addr) // 启动线程监听端口
	common.StopApplication()
	if !g.single {
		g.registerClient.Close()
	}
}

func (g *Gateway) ListenBack(uuid uint32, bytes []byte) []byte {
	if g.single {
		message := g.decoder.DecoderBytes(bytes)
		context := &router.Context{Message: message, SocketId: uuid}
		g.router.ExecuteMethod(context)
		if context.Data != nil {
			return context.Data
		}
		if context.Message == nil { // 没数据,底层socket对于空数据不返回
			return []byte{}
		}
		return context.Message.GetBytesResult()
	}
	// 此处调用远程方法
	ip, err := g.registerClient.GetIp()
	if err != nil {
		log.Println("注册中心错误:" + err.Error())
		return []byte{}
	}
	return g.client.InvokeRemoteRpc(ip.Addr(), bytes)
}
