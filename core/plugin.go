package core

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/router"
)

// Plugin 插件顶级接口
type Plugin interface {
	// GetId 唯一Id
	GetId() uint32
}

// GatewayPlugin 网关服插件必须实现的接口
type GatewayPlugin interface {
	// GetId 唯一Id
	GetId() uint32
	// SetService 设置逻辑服实例
	SetService(plugin *Gateway)
	// InvokeResult 回传信息
	InvokeResult([]byte) []byte
}

// ServicePlugin 逻辑服插件必须实现的接口
type ServicePlugin interface {
	// GetId 唯一Id
	GetId() uint32
	// SetService 设置逻辑服实例
	SetService(plugin *Service)
	// SetContext 设置,每次调用逻辑服都会执行!
	SetContext(ctx *router.Context)
}

type PluginService struct {
	// 插件
	pluginMap map[uint32]Plugin
}

// AddPlugin 添加插件
func (g *PluginService) AddPlugin(r Plugin) {
	if g.pluginMap == nil {
		g.pluginMap = make(map[uint32]Plugin)
	}
	if g.pluginMap[r.GetId()] != nil {
		panic("该插件:" + fmt.Sprint(r.GetId()) + "->Id已经存在不能重复添加!")
	}
	g.pluginMap[r.GetId()] = r
}

func (g *PluginService) UsePlugin(pluginId uint32, f func(r Plugin)) {
	v := g.pluginMap[pluginId]
	if v == nil {
		common.Logger().Println("Plugin: " + fmt.Sprint(pluginId) + " -> Id 不存在!")
		return
	}
	f(v)
}
