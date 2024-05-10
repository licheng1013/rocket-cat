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
func (d *DefaultRoom) JoinRoom(player IPlayer) {
	d.userMap.Store(player.UserId(), player)
}
func (d *DefaultRoom) QuitRoom(player IPlayer) {
	d.userMap.Delete(player.UserId())
}
