package room

import (
	"sync"
)

// DefaultRoom 默认房间实现,请继承此结构体,并重写方法,线程安全
type DefaultRoom struct {
	// 房间id
	RoomId int64
	// 用户Ids
	userMap sync.Map
	// 状态
	Status
	// 创建时间十位时间戳
	CreateTime int64
	// 房间管理器,注意如果此房间被删除了那么此值为nil
	manager *Manger
}

func (d *DefaultRoom) GetPlayerTotal() int {
	return len(d.GetUserIds())
}

func (d *DefaultRoom) ClearRoom() {
	d.userMap = sync.Map{}
}

func (d *DefaultRoom) GetPlayer(userId int64) IPlayer {
	if value, ok := d.userMap.Load(userId); ok {
		return value.(IPlayer)
	}
	return nil
}

func (d *DefaultRoom) GetId() int64 {
	return d.RoomId
}

func (d *DefaultRoom) GetState() Status {
	return d.Status
}

func (d *DefaultRoom) GetUserIds() (list []int64) {
	d.userMap.Range(func(key, value any) bool {
		list = append(list, key.(int64))
		return true
	})
	return
}

func (d *DefaultRoom) GetPlayers() []IPlayer {
	var userList []IPlayer
	d.userMap.Range(func(key, value any) bool {
		userList = append(userList, value.(IPlayer))
		return true
	})
	return userList
}

// 加入房间
func (d *DefaultRoom) JoinRoom(player IPlayer) bool {
	var join bool
	if value, ok := d.manager.roomIdOnRoom.Load(d.RoomId); ok {
		room := value.(IRoom)
		d.userMap.Store(player.UserId(), player)
		d.manager.userOnRoom.Store(player.UserId(), room)
		join = true
	}
	return join
}

/// 退出房间，如果房间内没有玩家则移除房间
func (d *DefaultRoom) QuitRoom(player IPlayer) {
	d.manager.userOnRoom.Delete(player.UserId())
	d.userMap.Delete(player.UserId())
	if d.GetPlayerTotal() == 0 {
		d.manager.RemoveRoom(d.RoomId)
		d.manager = nil // 释放
	}
}
