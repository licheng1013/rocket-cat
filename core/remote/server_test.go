package remote

import (
	"log"
	"net/http"
	"testing"
)

const addr = "192.168.101.10:10000"

func TestGrpcServer(t *testing.T) {
	server := GrpcServer{}
	server.CallbackResult(func(bytes []byte) []byte {
		log.Println("收到数据: ", string(bytes))
		return []byte("Hi")
	})
	server.ListenAddr(addr)
}

func TestGrpcClient(t *testing.T) {
	GrpcClientTest()
}

var grpcClient = GrpcClient{}

// 定义一个请求处理器
func helloHandler(w http.ResponseWriter, r *http.Request) {
	rpc := grpcClient.InvokeRemoteRpc(addr, []byte("HelloWorld"))
	_, _ = w.Write(rpc)
}

func GrpcClientTest() {
	// 将请求处理器注册到根路径上
	http.HandleFunc("/", helloHandler)
	// 启动一个 HTTP 服务器，监听 8080 端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
