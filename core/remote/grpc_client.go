package remote

import (
	"context"
	"github.com/io-game-go/protof"
	"google.golang.org/grpc"
	"log"
)

type GrpcClient struct {
	clientMap map[string]protof.RpcServiceClient
}

func (s *GrpcClient) InvokeRemoteRpc(addr string, bytes []byte) []byte {
	if len(addr) == 0 {
		log.Println("地址为空: " + addr)
		return []byte{}
	}
	if s.clientMap[addr] == nil {
		// 设置与服务器的连接
		socket, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("监听错误: %v", err)
		}
		s.clientMap[addr] = protof.NewRpcServiceClient(socket)
	}
	result, err := s.clientMap[addr].InvokeRemoteFunc(context.Background(), &protof.RpcInfo{Body: bytes})
	if err != nil {
		log.Fatalf("错误: %v", err)
	}
	return result.Body
}

func (s *GrpcClient) InvokeAllRemoteRpc(addrs []string, bytes []byte) {
	if len(addrs) == 0 {
		log.Println("找不到可用的服务端地址: ", addrs)
		return
	}
	for _, item := range addrs {
		if s.clientMap[item] == nil {
			s.InvokeRemoteRpc(item, bytes)
		}
	}
}
