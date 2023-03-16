package core

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/router"
)

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
	SetService(plugin *Service)
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
