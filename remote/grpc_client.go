package remote

import (
	"context"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/protof"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type GrpcClient struct {
	clientMap sync.Map //[string]protof.RpcServiceClient
}

func (s *GrpcClient) InvokeRemoteRpc(addr string, rpcInfo *protof.RpcInfo) []byte {
	if len(addr) == 0 {
		common.RocketLog.Println("地址为空: " + addr)
		return []byte{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, ok := s.clientMap.Load(addr)
	if !ok {
		// 设置与服务器的连接
		socket, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			common.RocketLog.Println("监听错误:" + err.Error())
		}
		value = protof.NewRpcServiceClient(socket)
		s.clientMap.Store(addr, value)
	}

	invoke, invokeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer invokeCancel()

	result, err := value.(protof.RpcServiceClient).InvokeRemoteFunc(invoke, rpcInfo)
	if err != nil {
		common.FileLogger().Println("请检查远程服务,远程错误:" + err.Error())
		s.clientMap.Delete(addr)
		return []byte{} //返回空则不返回给客户端
	}
	return result.Body
}

func (s *GrpcClient) InvokeAllRemoteRpc(addr []string, bytes []byte) {
	if len(addr) == 0 {
		common.RocketLog.Println("找不到可用的服务端地址: ", addr)
		return
	}
	for _, item := range addr {
		value, ok := s.clientMap.Load(item)
		if ok {
			_, _ = value.(protof.RpcServiceClient).InvokeRemoteFunc(context.Background(), &protof.RpcInfo{Body: bytes})
		}
	}
}
