package pkc

import (
	"fmt"
	"github.com/io-game-go/message"
	"github.com/io-game-go/registers"
)

// Rpc 远程调用接口,你可以随意实现自己的远程调用！
type Rpc interface {
	// Call 注册中心参数，路由，客户端的参数
	Call(requestUrl registers.RegisterInfo, info message.Message, rpcResult *RpcResult) error
	RpcListen(ip string, port uint64)
}

type RpcHandle interface {
	Invok(rpcInfo message.Message, rpcResulet *RpcResult) error
}

// Result 调用客户端返回的结果!
type Result struct {
}

func (r *Result) Invok(rpcInfo message.Message, rpcResulet *RpcResult) error {
	fmt.Println("收到信息: ", rpcInfo.GetMerge(), string(rpcInfo.GetBody()))
	rpcResulet.Result = []byte("Hello World")
	// TODO 这里是业务逻辑的处理！
	return nil
}

type RpcResult struct {
	Result []byte
	Error  string
}
