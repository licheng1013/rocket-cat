package remote

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/router"
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

// Plugin 插件
type Plugin interface {
	// InvokeResult 回传信息
	InvokeResult([]byte) []byte
	// GetId 唯一Id
	GetId() int32
}

// ServicePlugin 用于逻辑服的插件功能
type ServicePlugin interface {
	// InvokeResult 回传信息
	InvokeResult([]byte) []byte
	// GetId 唯一Id
	GetId() int32
	// SetService 设置逻辑服实例
	SetService(plugin *core.Service)
	// SetContext 设置,每次调用逻辑服都会执行!
	SetContext(ctx *router.Context)
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
