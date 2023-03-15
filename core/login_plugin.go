package core

import (
	"encoding/json"
	"log"
	"sync"
)

type LoginAction int8

const (
	Login = iota
	LogoutBySocketId
	LogoutByUserId
	ListSocketId
	ListUserId
)

type LoginPlugin struct {
	gateway     *Gateway
	userMap     sync.Map
	socketIdMap sync.Map
}

type LoginInterface interface {
	Login(userId int64, socketId uint32)
	LogoutBySocketId(socketId ...uint32)
	LogoutByUserId(userId ...int64)
	ListSocketId() (socketIds []uint32)
	ListUserId() (userIds []int64)
}

// LoginBody 登入数据
type LoginBody struct {
	LoginAction LoginAction
	UserId      []int64
	SocketId    []uint32
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

func (g *LoginPlugin) CallbackResult(bytes []byte) []byte {
	l := &LoginBody{}
	err := l.ToUnmarshal(bytes)
	if err != nil {
		log.Panicln("LoginBody -> 解析失败请检查或报告")
		return []byte{}
	}
	switch l.LoginAction {
	case Login:
        if len(l.UserId) == 1 && len(l.SocketId) == 1 {
        	g.Login(l.UserId[0],l.SocketId[0])
        }else{
			log.Println("LoginPlugin -> UserId或SocketId为空")
		}
		break
	case LogoutBySocketId:
		g.LogoutBySocketId(l.SocketId...)
		break
	case LogoutByUserId:
		g.LogoutByUserId(l.UserId...)
		break
	case ListSocketId:
		l.SocketId =  g.ListSocketId()
		break
	case ListUserId:
		l.UserId = g.ListUserId()
		break
	}
    marshal, err := l.ToMarshal()
    if err != nil {
		log.Panicln("LoginBody -> 解析失败请检查或报告")
		return []byte{}
    }
	return marshal
}

func (g *LoginPlugin) GetId() int32 {
	return 1
}

// Login 登入
func (g *LoginPlugin) Login(userId int64, socketId uint32) {
	g.userMap.Store(userId, socketId)
	g.socketIdMap.Store(socketId, userId)
}

// LogoutBySocketId 根据socketId退出
func (g *LoginPlugin) LogoutBySocketId(socketId ...uint32) {
	for _, id := range socketId {
		value, ok := g.socketIdMap.Load(id)
		if ok {
			g.userMap.Delete(value)
			g.socketIdMap.Delete(id)
		}
	}
}

// LogoutByUserId 根据用户id退出
func (g *LoginPlugin) LogoutByUserId(userId ...int64) {
	for _, id := range userId {
		value, ok := g.userMap.Load(id)
		if ok {
			g.userMap.Delete(id)
			g.socketIdMap.Delete(value)
		}
	}
}

// ListSocketId 获取所有客户端id
func (g *LoginPlugin) ListSocketId() (socketIds []uint32) {
	g.socketIdMap.Range(func(key, value any) bool {
		socketIds = append(socketIds, key.(uint32))
		return true
	})
	return
}

// ListUserId 获取所有用户id
func (g *LoginPlugin) ListUserId() (userIds []int64) {
	g.userMap.Range(func(key, value any) bool {
		userIds = append(userIds, key.(int64))
		return true
	})
	return
}
