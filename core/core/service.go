package core

import (
	"github.com/io-game-go/common"
	"github.com/io-game-go/decoder"
	"github.com/io-game-go/remote"
	"github.com/io-game-go/router"
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
	common.AssertNil(n.rpcServer, "The rpc server is not setup.")
	common.AssertNil(n.router, "The router is not setup.")
	common.AssertNil(n.decoder, "The decoder is not setup.")
	n.rpcServer.CallbackResult(func(bytes []byte) []byte { //这里回调数据，并进行内部处理
		message := n.decoder.DecoderBytes(bytes)
		context := router.Context{Message: message}
		return n.router.ExecuteMethod(context)
	})
	common.StopApplication()
	for _, item := range n.close {
		item() // Close Application
	}
}
