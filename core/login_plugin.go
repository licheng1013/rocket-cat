package core

import "sync"

type LoginPlugin struct {
	gateway     *Gateway
	userMap     sync.Map
	socketIdMap sync.Map
}

// Login 登入
func (g *LoginPlugin) Login(userId int64, sockdId uint32) {
	g.userMap.Store(userId, sockdId)
	g.socketIdMap.Store(sockdId, userId)
}

// LogoutBySocketId 根据socketId退出
func (g *LoginPlugin) LogoutBySocketId(sockdId ...uint32) {
	for _, id := range sockdId {
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
