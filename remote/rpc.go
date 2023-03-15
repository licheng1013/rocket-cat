package remote

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/protof"
)

// RpcClient 远程调用接口,你可以随意实现自己的远程调用！
type RpcClient interface {
	// InvokeAllRemoteRpc  只发送，不确认返回.
	InvokeAllRemoteRpc([]string, []byte)
	// InvokeRemoteRpc  Invoke remote method
	InvokeRemoteRpc(addr string, rpcInfo *protof.RpcInfo) []byte
}

type RpcServer interface {
	// CallbackResult 注册函数获取结果
	CallbackResult(func(in *protof.RpcInfo) []byte)
	// ListenAddr 监听地址
	ListenAddr(addr string)
}

// CallRpcInfo 回传信息
type CallRpcInfo interface {
	// CallbackResult 回传信息
	CallbackResult([]byte) []byte
	// GetId 唯一Id
	GetId() int32
}

// CallBody 回传信息
type CallBody struct {
	Id   int32
	Data []byte
}

// ToMarshal 转换为字节
func (b *CallBody) ToMarshal() (data []byte, err error) {
	data, err = json.Marshal(b)
	return
}

// ToUnmarshal 转换为对象
func (b *CallBody) ToUnmarshal(data []byte) (err error) {
	err = json.Unmarshal(data, b)
	return
}
