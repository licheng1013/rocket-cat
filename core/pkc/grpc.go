package pkc

import (
	"context"
	"core/message"
	"core/protof"
	"core/register"
	"flag"
	"fmt"
	"github.com/licheng1013/go-util/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

type Grpc struct {
	conn []*grpc.ClientConn
}

func (g *Grpc) Call(requestUrl register.RequestInfo, info message.Message, rpcResult *RpcResult) error {
	if len(g.conn) <= 100 { //TODO 设置最大连接数！
		conn, err := grpc.Dial(fmt.Sprintf("%v:%v", requestUrl.Ip, requestUrl.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		g.conn = append(g.conn, conn)
	}

	//defer g.conn.Close() //假设不关闭如何！
	c := protof.NewGrpcServiceClient(g.conn[common.RandomUtil.RandomRangeNum(len(g.conn))])
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	invoke, err := c.Invoke(ctx, info.(*protof.ProtoMessage))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	//log.Println("请求结果:" + string(invoke.Result))
	rpcResult.Result = invoke.GetResult()
	rpcResult.Error = invoke.Error
	return nil
}

// RpcListen 被调用逻辑处理
func (*Grpc) RpcListen(ip string, p uint64) {
	port := flag.Int("port", int(p), "The server port")
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d", ip, *port))
		if err != nil {
			log.Fatalf("监听失败: %v", err)
		}
		s := grpc.NewServer()
		protof.RegisterGrpcServiceServer(s, &server{})
		log.Printf("监听端口 %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	protof.UnimplementedGrpcServiceServer
}

func (s server) Invoke(ctx context.Context, in *protof.ProtoMessage) (*protof.GrpcResult, error) {
	return &protof.GrpcResult{Result: []byte("HelloWorld ok")}, nil
}
