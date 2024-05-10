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
	// RegisterServer Client
	registerClient registers.IRegister
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

func (g *Gateway) SetRegister(registerClient registers.IRegister) {
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

func (g *Gateway) Action(cmd, subCmd int64, method func(ctx *router.Context)) {
	g.router.Action(cmd, subCmd, method)
}

// 使用编码转换对象为数据
func (g *Gateway) data(data any) []byte {
	return g.decoder.Encode(data)
}

func (g *Gateway) ToRouterData(cmd, subCmd int64, data any) []byte {
	toData := g.data(data)
	return g.decoder.Data(cmd, subCmd, toData)
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
			go g.server.ListenAddr(g.registerClient.ClientInfo().Addr())
		}
	}
	g.socket.ListenBack(g.ListenBack)
	common.CatLog.Println("监听Socket:" + addr)
	go g.socket.ListenAddr(addr) // 启动线程监听端口
	common.StopApplication()
	if !g.single {
		g.registerClient.Close()
	}
}

func (g *Gateway) ListenBack(uuid uint32, bytes []byte) []byte {
	if g.single {
		message := g.decoder.Decoder(bytes)
		context := &router.Context{Message: message, SocketId: uuid}
		g.router.ExecuteMethod(context)
		if context.Data == nil {
			return []byte{}
		}
		context.Message.SetBody(context.Data)
		return context.Message.GetBytesResult()
	}
	// 此处调用远程方法
	ip, err := g.registerClient.GetIp()
	if err != nil {
		common.CatLog.Println("注册中心错误:" + err.Error())
		return []byte{}
	}
	remoteIp := g.socketIdIpMap[uuid]
	if remoteIp != "" {
		return g.client.InvokeRemoteRpc(remoteIp, &protof.RpcInfo{Body: bytes, SocketId: uuid, Ip: g.registerClient.ClientInfo().Addr()})
	}
	return g.client.InvokeRemoteRpc(ip.Addr(), &protof.RpcInfo{Body: bytes, SocketId: uuid, Ip: g.registerClient.ClientInfo().Addr()})
}

func (g *Gateway) SetServer(r remote.RpcServer) {
	g.server = r
	g.server.CallbackResult(g.CallbackResult)
}

// CallbackResult 给予远程端的回调方法
func (g *Gateway) CallbackResult(in *protof.RpcInfo) []byte {
	if g.pluginMap == nil || g.pluginMap[in.SocketId] == nil {
		common.CatLog.Println("插件不存在")
		return []byte{}
	}
	plugin := g.pluginMap[in.SocketId]
	return plugin.(GatewayPlugin).InvokeResult(in.GetBody())
}

// 发送给所有连接的客户端
func (g *Gateway) Push(result []byte) {
	g.socket.SendMessage(result)
}
