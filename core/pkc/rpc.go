package pkc

import (
	"core/register"
	"fmt"
)

// Rpc 远程调用接口,你可以随意实现自己的远程调用！
type Rpc interface {
	// Call 注册中心参数，路由，客户端的参数
	Call(requestUrl register.RequestInfo, info RequestInfo, rpcResult *RpcResult) error
}

// Result 结果
type Result struct {
}

func (r *Result) Invok(rpcInfo RequestInfo, rpcResulet *RpcResult) error {
	fmt.Println("收到信息: ", rpcInfo)
	rpcResulet.Result = []byte("Hello World")
	// TODO 这里是业务逻辑的处理！
	return nil
}
