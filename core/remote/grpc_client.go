package remote

import (
	"context"
	"github.com/io-game-go/protof"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type GrpcClient struct {
	clientMap sync.Map //[string]protof.RpcServiceClient
}

func (s *GrpcClient) InvokeRemoteRpc(addr string, bytes []byte) []byte {
	if len(addr) == 0 {
		log.Println("地址为空: " + addr)
		return []byte{}
	}
	value, ok := s.clientMap.Load(addr)

	if !ok {
		// 设置与服务器的连接
		socket, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Println("监听错误:" + err.Error())
		}
		value = protof.NewRpcServiceClient(socket)
		s.clientMap.Store(addr, value)
	}
	result, err := value.(protof.RpcServiceClient).InvokeRemoteFunc(context.Background(), &protof.RpcInfo{Body: bytes})
	if err != nil {
		log.Println("请检查远程服务,远程错误:" + err.Error())
		s.clientMap.Delete(addr)
		return []byte{} //返回空则不返回给客户端
	}
	return result.Body
}

func (s *GrpcClient) InvokeAllRemoteRpc(addrs []string, bytes []byte) {
	if len(addrs) == 0 {
		log.Println("找不到可用的服务端地址: ", addrs)
		return
	}
	for _, item := range addrs {
		value, ok := s.clientMap.Load(item)
		if ok {
			_, _ = value.(protof.RpcServiceClient).InvokeRemoteFunc(context.Background(), &protof.RpcInfo{Body: bytes})
		}
	}
}
