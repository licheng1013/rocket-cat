package remote

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/protof"
	"testing"
	"time"
)

func TestGrpcServer(t *testing.T) {
	const addr = "192.168.101.10:10002"
	server := GrpcServer{}
	server.CallbackResult(func(in *protof.RpcInfo) []byte {
		common.Logger().Println("收到数据: ", string(in.Body))
		return []byte("Hi")
	})
	go server.ListenAddr(addr)
	time.Sleep(time.Second)
	client := GrpcClient{}
	rpcResult := client.InvokeRemoteRpc(addr, protof.RpcBodyBuild([]byte("HelloWorld")))
	common.Logger().Println(string(rpcResult))
}
