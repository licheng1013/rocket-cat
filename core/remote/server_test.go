package remote

import (
	"context"
	"github.com/io-game-go/protof"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

const port = ":10000"

func TestGrpcServer(t *testing.T) {
	// 监听连接
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("监听: %v", err)
	}
	s := grpc.NewServer()
	protof.RegisterRpcServiceServer(s, &GrpcServer{})
	log.Printf("地址 %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("监听失败: %v", err)
	}
}

func TestGrpcClient(t *testing.T) {
	for i := 0; i < 50; i++ {
		go GrpcClient()
	}
	GrpcClient()
}

func GrpcClient() {
	// 设置与服务器的连接
	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("监听错误: %v", err)
	}
	defer conn.Close()
	c := protof.NewRpcServiceClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	for true {
		//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		r, err := c.InvokeRemoteFunc(context.Background(), &protof.RpcInfo{Body: []byte("HelloWorld")})
		if err != nil {
			log.Fatalf("错误: %v", err)
		}
		log.Println(r.String())
	}
	//defer cancel()
}
