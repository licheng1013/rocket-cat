package remote

import (
	"context"
	"github.com/io-game-go/protof"
)

type GrpcServer struct {
	protof.RpcServiceServer
}

func (s *GrpcServer) InvokeRemoteFunc(ctx context.Context, in *protof.RpcInfo) (*protof.RpcInfo, error) {
	in.Body = []byte("ok")
	return in, nil
}
