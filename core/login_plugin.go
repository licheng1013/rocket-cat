package core

import (
	"encoding/json"
	"sync"

	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/router"
)

// LoginPluginId 登入插件Id
const LoginPluginId = 1

const (
	Login = iota
	LogoutByUserId
	ListUserId
	SendAllUserMessage
	SendByUserIdMessage
	IsLogin
)

type LoginPlugin struct {
	gateway     *Gateway
	userMap     sync.Map
	socketIdMap sync.Map
}

func (g *LoginPlugin) OnClose(socketId uint32) {
	// 退出登入
	userId := g.ExistSocketId(socketId)
	if userId != 0 {
		common.RocketLog.Println("用户断开 -> ", userId)
		g.userMap.Delete(userId)
	}
}

// ExistSocketId 根据socketId判断是否登入
func (g *LoginPlugin) ExistSocketId(socketId uint32) int64 {
	value, ok := g.socketIdMap.Load(socketId)
	if ok {
		return value.(int64)
	}
	return 0
}

// SendAllUserMessage 广播所有登入用户
func (g *LoginPlugin) SendAllUserMessage(data []byte) {
	var socketId []uint32
	g.userMap.Range(func(key, value any) bool {
		socketId = append(socketId, value.(uint32))
		return true
	})
	g.gateway.socket.SendSelectMessage(data, socketId...)
}

// SendByUserIdMessage 根据用户id进行广播
func (g *LoginPlugin) SendByUserIdMessage(data []byte, userIds ...int64) {
	var socketId []uint32
	for _, userId := range userIds {
		value, ok := g.userMap.Load(userId)
		if ok {
			socketId = append(socketId, value.(uint32))
		}
	}
	g.gateway.socket.SendSelectMessage(data, socketId...)
}

func (g *LoginPlugin) SetService(plugin *Gateway) {
	g.gateway = plugin
}

type LoginInterface interface {
	Login(userId int64, socketId uint32) bool
	LogoutByUserId(userId int64) bool
	ListUserId() (userIds []int64)
	SendAllUserMessage(data []byte)
	SendByUserIdMessage(data []byte, userIds ...int64)
	IsLogin(userId int64) bool
}

// LoginBody 登入数据
type LoginBody struct {
	// 动作逻辑服需要调用网关服的那些方法
	Action uint8
	// 所有登入id
	UserIds []int64
	// 用户Id
	UserId int64
	// 连接建立时的id
	SocketId uint32
	// 登入或退出状态
	State bool
	// 广播数据
	Data []byte
}

// ToMarshal 转换为字节
func (b *LoginBody) ToMarshal() (data []byte) {
	data, err := json.Marshal(b)
	if err != nil {
		common.RocketLog.Println("json转换失败: " + err.Error())
	}
	if data == nil { //返回空
		return []byte{}
	}
	return data
}

// ToUnmarshal 转换为对象
func (b *LoginBody) ToUnmarshal(data []byte) {
	err := json.Unmarshal(data, b)
	if err != nil {
		common.RocketLog.Println("json解析失败:" + err.Error())
	}
	return
}

func (g *LoginPlugin) InvokeResult(bytes []byte) []byte {
	l := &LoginBody{}
	l.ToUnmarshal(bytes)
	switch l.Action {
	case Login:
		if l.UserId != 0 && l.SocketId != 0 {
			l.State = g.Login(l.UserId, l.SocketId)
		} else {
			common.RocketLog.Println("LoginPlugin -> UserId或SocketId为空")
		}
		break
	case LogoutByUserId:
		l.State = g.LogoutByUserId(l.UserId)
		break
	case ListUserId:
		l.UserIds = g.ListUserId()
		break
	case SendAllUserMessage:
		g.SendAllUserMessage(l.Data)
		break
	case SendByUserIdMessage:
		g.SendByUserIdMessage(l.Data, l.UserIds...)
		break
	case IsLogin:
		l.State = g.IsLogin(l.UserId)
		break
	}
	return l.ToMarshal()
}

func (g *LoginPlugin) GetId() uint32 {
	return LoginPluginId
}

// GetSocketIdByUserId 获取连接id
// 逻辑服可能不需要此方法，暂时不帮逻辑服实现
func (g *LoginPlugin) GetSocketIdByUserId(userId int64) uint32 {
	value, ok := g.userMap.Load(userId)
	if ok {
		return value.(uint32)
	}
	return 0
}

// IsLogin 是否登入
func (g *LoginPlugin) IsLogin(userId int64) bool {
	_, ok := g.userMap.Load(userId)
	return ok
}

// Login 登入,已存在则为false
func (g *LoginPlugin) Login(userId int64, socketId uint32) bool {
	_, ok := g.userMap.Load(userId) // 第一次肯定是空即false,否则就是已登入
	if !ok {
		g.userMap.Store(userId, socketId)
		g.socketIdMap.Store(socketId, userId)
	}
	return !ok
}

// LogoutByUserId 根据用户id退出
func (g *LoginPlugin) LogoutByUserId(userId int64) bool {
	value, ok := g.userMap.Load(userId)
	if ok {
		g.userMap.Delete(userId)
		g.socketIdMap.Delete(value)
	}
	return ok
}

// ListUserId 获取所有用户id
func (g *LoginPlugin) ListUserId() (userIds []int64) {
	g.userMap.Range(func(key, value any) bool {
		userIds = append(userIds, key.(int64))
		return true
	})
	return
}

type LoginPluginService struct {
	service *Service
	ctx     *router.Context
}

func (l *LoginPluginService) SendAllUserMessage(data []byte) {
	body := LoginBody{Action: SendAllUserMessage, Data: data}
	_, _ = l.service.SendGatewayMessage(&protof.RpcInfo{SocketId: LoginPluginId, Body: body.ToMarshal()})
}

func (l *LoginPluginService) SendByUserIdMessage(data []byte, userIds ...int64) {
	body := LoginBody{Action: SendByUserIdMessage, Data: data, UserIds: userIds}
	_, _ = l.service.SendGatewayMessage(&protof.RpcInfo{SocketId: LoginPluginId, Body: body.ToMarshal()})
}

// IsLogin 是否登入
func (l *LoginPluginService) IsLogin(userId int64) bool {
	body := LoginBody{UserId: userId, Action: IsLogin}
	message, err := l.service.SendGatewayMessage(&protof.RpcInfo{SocketId: LoginPluginId, Body: body.ToMarshal()})
	if err != nil {
		return false
	}
	for _, bytes := range message { //遍历所有结果知道有一个true那即为退出成功!
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

// Login 登入肯定是，登入当前连接，难道你从a网关登入到b网关去吗。。。
func (l *LoginPluginService) Login(userId int64, socketId uint32) bool {
	body := LoginBody{Action: Login, UserId: userId, SocketId: socketId}
	rpc := l.service.rpcClient.InvokeRemoteRpc(l.ctx.RpcIp, &protof.RpcInfo{SocketId: LoginPluginId, Body: body.ToMarshal()})
	if len(rpc) == 0 { // 此处数据为空必然是网关服出现问题
		return false
	}
	body.ToUnmarshal(rpc)
	return body.State
}

// LogoutByUserId 根据用户id退出,这里需要广播所有网关进行退出登入操作,因为逻辑服并不知道用户登入在那个网关服
func (l *LoginPluginService) LogoutByUserId(userId int64) bool {
	body := LoginBody{UserId: userId, Action: LogoutByUserId}
	message, err := l.service.SendGatewayMessage(&protof.RpcInfo{SocketId: LoginPluginId, Body: body.ToMarshal()})
	if err != nil {
		return false
	}
	for _, bytes := range message { //遍历所有结果知道有一个true那即为退出成功!
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

func (l *LoginPluginService) ListUserId() (userIds []int64) {
	body := LoginBody{Action: ListUserId}
	message, err := l.service.SendGatewayMessage(&protof.RpcInfo{SocketId: LoginPluginId, Body: body.ToMarshal()})
	if err != nil {
		return []int64{}
	}
	for _, bytes := range message {
		if len(bytes) == 0 {
			continue
		}
		body.ToUnmarshal(bytes)
		userIds = append(userIds, body.UserIds...)
	}
	return
}

func (l *LoginPluginService) SetContext(ctx *router.Context) {
	l.ctx = ctx
}

func (l *LoginPluginService) GetId() uint32 {
	return LoginPluginId
}

func (l *LoginPluginService) SetService(service *Service) {
	l.service = service
}
