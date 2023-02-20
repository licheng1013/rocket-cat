package core

import (
	"github.com/io-game-go/common"
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
}

func (n *Service) AddClose(close func()) {
	n.close = append(n.close, close)
}

func (n *Service) SetRpcServer(rpcServer remote.RpcServer) {
	n.rpcServer = rpcServer
}

func (n *Service) Router() router.Router {
	return n.router
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
	common.StopApplication()
	for _, item := range n.close {
		item() // Close Application
	}
}
