package remote

import (
	"github.com/licheng1013/io-game-go/protof"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const addr = "192.168.101.10:10001"

// 测试请求 http -> grpc
func TestGrpcClient(t *testing.T) {
	server := GrpcServer{}
	server.CallbackResult(func(in *protof.RpcInfo) []byte {
		log.Println("收到数据: ", string(in.Body))
		return []byte("Hi")
	})
	go server.ListenAddr(addr)
	go GrpcClientTest() //启动两个后访问 http://localhost:8080/
	time.Sleep(time.Second / 2)
	// 创建一个客户端对象
	client := &http.Client{}
	// 创建一个请求对象
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("返回结果:", string(body))
}

var grpcClient = GrpcClient{}

func GrpcClientTest() {
	// 将请求处理器注册到根路径上
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rpc := grpcClient.InvokeRemoteRpc(addr, protof.RpcBodyBuild([]byte("HelloWorld")))
		_, _ = writer.Write(rpc)
	})
	// 启动一个 HTTP 服务器，监听 8080 端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
