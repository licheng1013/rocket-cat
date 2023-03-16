package core

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/protof"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"sync"
)

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
	LoginAction LoginAction
	UserIds     []int64
	UserId      int64
	SocketId    uint32
}

// ToMarshal 转换为字节
func (b *LoginBody) ToMarshal() (data []byte, err error) {
	data, err = json.Marshal(b)
	return
}

// ToUnmarshal 转换为对象
func (b *LoginBody) ToUnmarshal(data []byte) (err error) {
	err = json.Unmarshal(data, b)
	return
}

func (g *LoginPlugin) InvokeResult(bytes []byte) []byte {
	l := &LoginBody{}
	err := l.ToUnmarshal(bytes)
	if err != nil {
		log.Panicln("LoginBody -> 解析失败请检查或报告")
		return []byte{}
	}
	switch l.LoginAction {
	case Login:
		if l.UserId != 0 && l.SocketId != 0 {
			g.Login(l.UserId, l.SocketId)
		} else {
			log.Println("LoginPlugin -> UserId或SocketId为空")
		}
		break
	case LogoutByUserId:
		g.LogoutByUserId(l.UserId)
		break
	case ListUserId:
		l.UserIds = g.ListUserId()
		break
	}
	marshal, err := l.ToMarshal()
	if err != nil {
		log.Panicln("LoginBody -> 解析失败请检查或报告")
		return []byte{}
	}
	return marshal
}

const pluginId = 1

func (g *LoginPlugin) GetId() int32 {
	return pluginId
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

func (l *LoginPluginService) Login(userId int64, socketId uint32) bool {
	l.service.rpcClient.InvokeRemoteRpc(l.ctx.RpcIp,&protof.RpcInfo{})
	return false
}

func (l *LoginPluginService) LogoutByUserId(userId int64) bool {
	//TODO implement me
	panic("implement me")
}

func (l *LoginPluginService) ListUserId() (userIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (l *LoginPluginService) SetContext(ctx *router.Context) {
	l.ctx = ctx
}

func (l *LoginPluginService) InvokeResult(bytes []byte) []byte {
	data, _ := l.service.SendGatewayMessage(bytes)
	return data[0]
}

func (l *LoginPluginService) GetId() int32 {
	return pluginId
}

func (l *LoginPluginService) SetService(service *Service) {
	l.service = service
}
