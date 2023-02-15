package remote

import (
	"context"
	"github.com/io-game-go/protof"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
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
	GrpcClient()
}

var conn *grpc.ClientConn
var client protof.RpcServiceClient

// 定义一个请求处理器
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if conn == nil {
		// 设置与服务器的连接
		c, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("监听错误: %v", err)
		}
		conn = c
		client = protof.NewRpcServiceClient(conn)
	}
	v, err := client.InvokeRemoteFunc(context.Background(), &protof.RpcInfo{Body: []byte("HelloWorld")})
	if err != nil {
		log.Fatalf("错误: %v", err)
	}
	log.Println(v.String())
	_, _ = w.Write([]byte("Hello, world!"))
}

func GrpcClient() {
	// 将请求处理器注册到根路径上
	http.HandleFunc("/", helloHandler)
	// 启动一个 HTTP 服务器，监听 8080 端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
