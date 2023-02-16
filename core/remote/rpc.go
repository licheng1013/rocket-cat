package remote

// Rpc 远程调用接口,你可以随意实现自己的远程调用！
type Rpc interface {
	// InvokeAllRemoteRpc  发送，
	InvokeAllRemoteRpc([]string, []byte)
	// InvokeRemoteRpc  发送，
	InvokeRemoteRpc(string, []byte) []byte
}
