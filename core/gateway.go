package core

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/version"
	"log"

	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
	router "github.com/licheng1013/rocket-cat/router"
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
	// Rpc 客户端
	client remote.RpcClient
	// Rpc 服务端
	server remote.RpcServer
	// Register Client
	registerClient registers.Register
	// 插件系统
	PluginService
	// 绑定 socketId 和 ip
	socketIdIpMap map[uint32]string
	// 关闭钩子
	closeHook []connect.SocketClose
}

func (g *Gateway) Socket() connect.Socket {
	return g.socket
}

func (g *Gateway) OnClose(socketId uint32) {
	for _, item := range g.closeHook {
		item.OnClose(socketId)
	}
}

func (g *Gateway) SetSocket(socket connect.Socket) {
	g.socket = socket
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

func DefaultGateway() *Gateway {
	g := &Gateway{}
	g.SetSingle(true)
	g.router = &router.DefaultRouter{}
	g.AddPlugin(&LoginPlugin{})
	g.AddPlugin(&BindPlugin{})
	g.SetDecoder(decoder.JsonDecoder{})
	g.SetSocket(&connect.WebSocket{})
	return g
}

func (g *Gateway) Start(addr string) {
	version.StartLogo()
	g.router.AddProxy(&router.ErrProxy{}) // 添加错误代理
	// 插件初始化
	for _, item := range g.PluginService.pluginMap {
		switch item.(type) {
		case connect.SocketClose:
			g.closeHook = append(g.closeHook, item.(connect.SocketClose))
		}
		switch item.(type) {
		case GatewayPlugin:
			item.(GatewayPlugin).SetService(g)
			break
		default:
			panic(fmt.Sprintf("此插件: %v 并没有实现 GatewayPlugin 接口", item.GetId()))
		}
	}
	// 设置断开钩子
	g.socket.OnClose(g.OnClose)
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	common.AssertNil(g.socket, "没有设置链接协议")
	if g.single {
		common.AssertNil(g.decoder, "没有设置编码器: decoder")
	} else {
		common.AssertNil(g.client, "没有设置远程客户端!")
		common.AssertNil(g.registerClient, "没有设置注册客户端!")
		if g.server != nil { // Rpc 服务端
			g.server.CallbackResult(func(in *protof.RpcInfo) []byte {
				// TODO 这是从 -> 逻辑服发送过来的消息
				return nil
			})
			go g.server.ListenAddr(g.registerClient.RegisterInfo().Addr())
		}
	}
	g.socket.ListenBack(g.ListenBack)
	common.Logger().Println("监听Socket:" + addr)
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
		common.Logger().Println("注册中心错误:" + err.Error())
		return []byte{}
	}
	remoteIp := g.socketIdIpMap[uuid]
	if remoteIp != "" {
		return g.client.InvokeRemoteRpc(remoteIp, &protof.RpcInfo{Body: bytes, SocketId: uuid, Ip: g.registerClient.RegisterInfo().Addr()})
	}
	return g.client.InvokeRemoteRpc(ip.Addr(), &protof.RpcInfo{Body: bytes, SocketId: uuid, Ip: g.registerClient.RegisterInfo().Addr()})
}

func (g *Gateway) SetServer(r remote.RpcServer) {
	g.server = r
	g.server.CallbackResult(g.CallbackResult)
}

// CallbackResult 给予远程端的回调方法
func (g *Gateway) CallbackResult(in *protof.RpcInfo) []byte {
	if g.pluginMap == nil || g.pluginMap[in.SocketId] == nil {
		common.Logger().Println("插件不存在")
		return []byte{}
	}
	plugin := g.pluginMap[in.SocketId]
	return plugin.(GatewayPlugin).InvokeResult(in.GetBody())
}

// SendMessage 广播所有
func (g *Gateway) SendMessage(result []byte) {
	g.socket.SendMessage(result)
}
