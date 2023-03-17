package core

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"sync"
)

// LoginPluginId 登入插件Id
const LoginPluginId = 1

// LoginAction  插件内部路由
type LoginAction int8

const (
	Login = iota
	LogoutByUserId
	ListUserId
)

type LoginPlugin struct {
	gateway     *Gateway
	userMap     sync.Map
	socketIdMap sync.Map
}

type LoginInterface interface {
	Login(userId int64, socketId uint32) bool
	LogoutByUserId(userId int64) bool
	ListUserId() (userIds []int64)
}

// LoginBody 登入数据
type LoginBody struct {
	// 动作逻辑服需要调用网关服的那些方法
	Action LoginAction
	// 所有登入id
	UserIds []int64
	// 用户Id
	UserId int64
	// 连接建立时的id
	SocketId uint32
	// 登入或退出状态
	State bool
}

// ToMarshal 转换为字节
func (b *LoginBody) ToMarshal() (data []byte) {
	data, err := json.Marshal(b)
	if err != nil {
		log.Println("json转换失败: " + err.Error())
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
		log.Println("json解析失败:" + err.Error())
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
			log.Println("LoginPlugin -> UserId或SocketId为空")
		}
		break
	case LogoutByUserId:
		l.State = g.LogoutByUserId(l.UserId)
		break
	case ListUserId:
		l.UserIds = g.ListUserId()
		break
	}
	return l.ToMarshal()
}

func (g *LoginPlugin) GetId() uint32 {
	return LoginPluginId
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
	message, err := l.service.SendGatewayMessage(body.ToMarshal())
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
	message, err := l.service.SendGatewayMessage(body.ToMarshal())
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
