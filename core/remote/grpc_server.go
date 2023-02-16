package remote

import (
	"context"
	"github.com/io-game-go/protof"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcServer struct {
	protof.RpcServiceServer
	callbackFunc func([]byte) []byte
}

func (s *GrpcServer) CallbackResult(f func([]byte) []byte) {
	s.callbackFunc = f
}

func (s *GrpcServer) ListenAddr(addr string) {
	// 监听连接
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("监听: %v", err)
	}
	v := grpc.NewServer()
	protof.RegisterRpcServiceServer(v, s)
	log.Printf("地址 %v", lis.Addr())
	if err := v.Serve(lis); err != nil {
		log.Fatalf("监听失败: %v", err)
	}
}

// InvokeRemoteFunc 此处由Grpc客户端调用
func (s *GrpcServer) InvokeRemoteFunc(ctx context.Context, in *protof.RpcInfo) (*protof.RpcInfo, error) {
	if s.callbackFunc == nil {
		panic("没有注册回调方法")
	}
	in.Body = s.callbackFunc(in.Body)
	return in, nil
}
