package core

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/router"
)

// BindPluginId 绑定插件Id
const BindPluginId = 2

const (
	Bind = iota
	UnBind
)

// BindInterface 单机模式中不需要此插件
type BindInterface interface {
	// Bind 多个玩家绑定到ip
	Bind(boyd BindBody) bool
	// UnBind 接触绑定
	UnBind(body BindBody) bool
}

// BindPlugin 绑定某个逻辑或解绑某个逻辑服务
type BindPlugin struct {
	gateway *Gateway
}

func (b *BindPlugin) Bind(body BindBody) bool {
	plugin := b.gateway.GetPlugin(LoginPluginId) //依赖于登入插件使用
	if plugin == nil {
		return false
	}
	loginPlugin := plugin.(*LoginPlugin)
	for _, item := range body.UserIds {
		if loginPlugin.IsLogin(item) { // 只有登入过才允许绑定插件
			id := loginPlugin.GetSocketIdByUserId(item)
			b.gateway.socketIdIpMap[id] = body.Ip
		}
	}
	return true
}

func (b *BindPlugin) UnBind(body BindBody) bool {
	plugin := b.gateway.GetPlugin(LoginPluginId) //依赖于登入插件使用
	if plugin == nil {
		return false
	}
	loginPlugin := plugin.(*LoginPlugin)
	for _, item := range body.UserIds {
		delete(b.gateway.socketIdIpMap, loginPlugin.GetSocketIdByUserId(item))
	}
	return true
}

func (b *BindPlugin) GetId() uint32 {
	return BindPluginId
}

func (b *BindPlugin) SetService(core *Gateway) {
	b.gateway = core
}

func (b *BindPlugin) InvokeResult(bytes []byte) []byte {
	body := &BindBody{}
	body.ToUnmarshal(bytes)
	switch body.Action {
	case Bind:
		body.State = b.Bind(*body)
		break
	case UnBind:
		body.State = b.UnBind(*body)
		break
	}
	return body.ToMarshal()
}

type BindBody struct {
	// 动作逻辑服需要调用网关服的那些方法
	Action uint8
	// 所有登入id
	UserIds []int64
	// 绑定Ip
	Ip string
	// 操作状态
	State bool
}

// ToMarshal 转换为字节
func (b *BindBody) ToMarshal() (data []byte) {
	data, err := json.Marshal(b)
	if err != nil {
		common.Logger().Println("json转换失败: " + err.Error())
	}
	if data == nil { //返回空
		return []byte{}
	}
	return data
}

// ToUnmarshal 转换为对象
func (b *BindBody) ToUnmarshal(data []byte) {
	err := json.Unmarshal(data, b)
	if err != nil {
		common.Logger().Println("json解析失败:" + err.Error())
	}
	return
}

type BindPluginService struct {
	service *Service
	ctx     *router.Context
}

func (b *BindPluginService) Bind(body BindBody) bool {
	body.Action = Bind
	message, err := b.service.SendGatewayMessage(&protof.RpcInfo{SocketId: BindPluginId, Body: body.ToMarshal()})
	if err != nil {
		return false
	}
	for _, bytes := range message {
		if len(bytes) == 0 {
			continue
		}
		body.ToUnmarshal(bytes)
		if body.State {
			return true
		}
	}
	return false
}

func (b *BindPluginService) UnBind(body BindBody) bool {
	body.Action = UnBind
	message, err := b.service.SendGatewayMessage(&protof.RpcInfo{SocketId: BindPluginId, Body: body.ToMarshal()})
	if err != nil {
		return false
	}
	for _, bytes := range message {
		if len(bytes) == 0 {
			continue
		}
		body.ToUnmarshal(bytes)
		if body.State {
			return true
		}
	}
	return false
}

func (b *BindPluginService) GetId() uint32 {
	return BindPluginId
}

func (b *BindPluginService) SetService(service *Service) {
	b.service = service
}

// SetContext 每次调用插件都会处理ctx
func (b *BindPluginService) SetContext(ctx *router.Context) {
	b.ctx = ctx
}
