package remote

// RpcClient 远程调用接口,你可以随意实现自己的远程调用！
type RpcClient interface {
	// InvokeAllRemoteRpc  只发送，不确认返回.
	InvokeAllRemoteRpc([]string, []byte)
	// InvokeRemoteRpc  发送，
	InvokeRemoteRpc(string, []byte) []byte
}

type RpcServer interface {
	// CallbackResult 注册函数获取结果
	CallbackResult(func([]byte) []byte)
	// ListenAddr 监听地址
	ListenAddr(addr string)
}
