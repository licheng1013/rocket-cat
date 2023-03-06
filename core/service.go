package core

import (
	"github.com/licheng1013/io-game-go/common"
	"github.com/licheng1013/io-game-go/decoder"
	"github.com/licheng1013/io-game-go/protof"
	"github.com/licheng1013/io-game-go/registers"
	"github.com/licheng1013/io-game-go/remote"
	"github.com/licheng1013/io-game-go/router"
	"log"
)

// Service 新手请不需要之间使用而是 NewService 进行获取对象
type Service struct {
	// 路由
	router router.Router
	// rpc 监听
	rpcServer remote.RpcServer
	// 关机钩子
	close []func()
	// 编码器
	decoder decoder.Decoder
	// 注册中心
	register registers.Register
}

func (n *Service) Router() router.Router {
	return n.router
}

func (n *Service) SetDecoder(decoder decoder.Decoder) {
	n.decoder = decoder
}

func (n *Service) AddClose(close func()) {
	n.close = append(n.close, close)
}
func (n *Service) SetRpcServer(rpcServer remote.RpcServer) {
	n.rpcServer = rpcServer
}
func (n *Service) SetRouter(router router.Router) {
	n.router = router
}
func NewService() *Service {
	service := &Service{}
	service.SetRouter(&router.DefaultRouter{})
	return service
}

func (n *Service) Start() {
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	common.AssertNil(n.rpcServer, "Rpc服务没有设置.")
	common.AssertNil(n.router, "路由没有设置.")
	common.AssertNil(n.decoder, "编码器没有设置.")
	common.AssertNil(n.register, "注册中心没有设置.")
	n.rpcServer.CallbackResult(func(bytes []byte) []byte { //这里回调数据，并进行内部处理
		p := &protof.RpcBody{}
		protof.RpcBodyUnmarshal(bytes, p)
		log.Println(p.Body)
		message := n.decoder.DecoderBytes(p.Body)
		context := &router.Context{Message: message, SocketId: p.SocketId}
		context.RpcServer = n.rpcServer
		n.router.ExecuteMethod(context)
		if context.Data != nil {
			return context.Data
		}
		if context.Message == nil {
			return []byte{}
		}
		return context.Message.GetBytesResult()
	})
	n.close = append(n.close, n.register.Close)
	addr := n.register.RegisterInfo().Addr()
	go n.rpcServer.ListenAddr(addr)
	common.StopApplication()
	for _, item := range n.close {
		item() // Close Application
	}
}

func (n *Service) SetRegister(register registers.Register) {
	n.register = register
}
