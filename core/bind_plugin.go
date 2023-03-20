package core

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
	//TODO implement me
	panic("implement me")
}


type BindBody struct {
	// 动作逻辑服需要调用网关服的那些方法
	Action uint8
	// 所有登入id
	UserIds []int64
	// 绑定Ip
	Ip string
}