package remote

import (
	"context"
	"fmt"
	"github.com/io-game-go/protof"
	"sync"
	"time"
)

type GrpcServer struct {
	protof.RpcServiceServer
}

var start = time.Now().UnixMilli()
var count int64
var lock sync.Mutex

func (s *GrpcServer) InvokeRemoteFunc(ctx context.Context, in *protof.RpcInfo) (*protof.RpcInfo, error) {
	end := time.Now().UnixMilli()
	lock.Lock()
	count++
	if end-start > 1000 {
		fmt.Println("1s请求数量:", count)
		count = 0
		start = end
	}
	lock.Unlock()
	//log.Println("收到:", in.String())
	in.Body = []byte("ok")
	return in, nil
}
