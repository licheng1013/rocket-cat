package remote

import "github.com/licheng1013/io-game-go/registers"

// RpcClient 远程调用接口,你可以随意实现自己的远程调用！
type RpcClient interface {
	// InvokeAllRemoteRpc  只发送，不确认返回.
	InvokeAllRemoteRpc([]string, []byte)
	// InvokeRemoteRpc  Invoke remote method
	InvokeRemoteRpc(addr string, bytes []byte) []byte
}

type RpcServer interface {
	// CallbackResult 注册函数获取结果
	CallbackResult(func([]byte) []byte)
	// ListenAddr 监听地址
	ListenAddr(addr registers.RegisterInfo)
	// SetRegister 设置注册中心
	SetRegister(register registers.Register)
}
