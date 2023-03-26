package core

import (
	"errors"
	"fmt"
	"github.com/licheng1013/rocket-cat/version"
	"log"
	"time"

	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
	"github.com/licheng1013/rocket-cat/router"
)

// Service 新手请不需要之间使用而是 NewService 进行获取对象
type Service struct {
	// 路由
	router router.Router
	// rpc 监听
	rpcServer remote.RpcServer
	// rpc 客户端
	rpcClient remote.RpcClient
	// 关机钩子
	close []func()
	// 编码器
	decoder decoder.Decoder
	// 注册中心
	register registers.Register
	// 线程池,用于请求多逻辑服事使用
	Pool *common.Pool
	// 插件系统
	PluginService
}

// SendGatewayMessage 广播消息路由 -> 所有网关服
func (n *Service) SendGatewayMessage(info *protof.RpcInfo) (result [][]byte, err error) {
	return n.sendMessageByServiceName(n.register.RegisterInfo().RemoteName, info)
}

// SendServiceMessage 广播消息路由 -> 所有逻辑服
func (n *Service) SendServiceMessage(bytes []byte) (result [][]byte, err error) {
	return n.sendMessageByServiceName(n.register.RegisterInfo().ServiceName, protof.RpcBodyBuild(bytes))
}

// sendMessageByServiceName  广播消息路由
func (n *Service) sendMessageByServiceName(serviceName string, rpcInfo *protof.RpcInfo) (result [][]byte, err error) {
	ips, err := n.register.ListIp(serviceName)
	if len(ips) == 0 {
		common.Logger().Println("注册中心暂无可用的服务!")
		return [][]byte{}, nil
	}
	if err != nil {
		//common.Logger().Println("注册中心错误 -> " + err.Error())
		return nil, err
	}
	channel := make(chan []byte)
	for _, item := range ips {
		n.Pool.AddTask(func() {
			channel <- n.rpcClient.InvokeRemoteRpc(item.Addr(), rpcInfo)
		})
	}
	for i := 0; i < len(ips); i++ {
		select {
		case b := <-channel:
			if len(b) != 0 {
				result = append(result, b)
			}
		case <-time.After(2 * time.Second):
			err = errors.New("有逻辑服连接超时")
		}
	}
	return
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

func DefaultService() *Service {
	service := &Service{}
	service.SetRouter(&router.DefaultRouter{})
	service.AddPlugin(&LoginPluginService{})
	service.AddPlugin(&BindPluginService{})
	return service
}

// CallbackResult 回调数据
func (n *Service) CallbackResult(in *protof.RpcInfo) []byte {
	message := n.decoder.DecoderBytes(in.Body)
	context := &router.Context{Message: message, SocketId: in.SocketId, RpcIp: in.Ip}
	context.RpcServer = n.rpcServer
	for _, plugin := range n.pluginMap {
		switch plugin.(type) {
		case ServicePlugin:
			plugin.(ServicePlugin).SetContext(context)
			break
		default:
			log.Printf("此插件Id: %v 没有实现 -> ServicePlugin 接口所以将不会执行\n", plugin.GetId())
		}
	}
	n.router.ExecuteMethod(context)
	if context.Data != nil {
		return context.Data
	}
	if context.Message == nil {
		return []byte{}
	}
	return context.Message.GetBytesResult()
}

func (n *Service) Start() {
	version.StartLogo()
	// 插件初始化
	for _, item := range n.PluginService.pluginMap {
		switch item.(type) {
		case ServicePlugin:
			item.(ServicePlugin).SetService(n)
			break
		default:
			panic(fmt.Sprintf("此插件: %v 并没有实现 ServicePlugin 接口", item.GetId()))
		}
	}

	log.SetFlags(log.LstdFlags + log.Lshortfile)
	common.AssertNil(n.rpcServer, "Rpc服务没有设置.")
	common.AssertNil(n.router, "路由没有设置.")
	common.AssertNil(n.decoder, "编码器没有设置.")
	common.AssertNil(n.register, "注册中心没有设置.")
	if n.Pool == nil {
		n.Pool = common.NewPool()
	}
	n.rpcServer.CallbackResult(n.CallbackResult)
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

func (n *Service) SetRpcClient(r *remote.GrpcClient) {
	n.rpcClient = r
}
