package core

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/common"
)

// BindPluginId 绑定插件Id
const BindPluginId = 2

const (
	Bind = iota
	UnBind
)

// BindInterface 单机模式中不需要此插件
type BindInterface interface {
	Bind(userId []int64, ip string) bool
	UnBind()
}


// BindPlugin 绑定某个逻辑或解绑某个逻辑服务
type BindPlugin struct {
	gateway     *Gateway
}

func (b *BindPlugin) GetId() uint32 {
	return BindPluginId
}

func (b *BindPlugin) SetService(plugin *Gateway) {
	b.gateway = plugin
}

func (b *BindPlugin) InvokeResult(bytes []byte) []byte {
	body := &BindBody{}
	body.ToUnmarshal(bytes)
	return []byte{}
}


type BindBody struct {
	// 动作逻辑服需要调用网关服的那些方法
	Action uint8
	// 所有登入id
	UserIds []int64
	// 绑定Ip
	Ip string
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