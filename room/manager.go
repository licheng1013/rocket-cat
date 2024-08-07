package room

import (
	"sync"
)

type Manger struct {
	// 用户id - 房间id
	userOnRoom sync.Map
	// roomId - 房间具体房间
	roomIdOnRoom sync.Map
}

func NewManger() *Manger {
	return &Manger{}
}

// GetUniqueRoomId 创建房间
func (m *Manger) GetUniqueRoomId() int64 {
	/// 获取最大的房间id+1返回
	var maxRoodId int64
	m.roomIdOnRoom.Range(func(key, value any) bool {
		if key.(int64) > maxRoodId {
			maxRoodId = key.(int64)
		}
		return true
	})
	return maxRoodId + 1
}

func (m *Manger) GetByUserId(userId int64) IRoom {
	value, ok := m.userOnRoom.Load(userId)
	if ok {
		return value.(IRoom)
	}
	return nil
}

// AddRoom 添加房间, 无需手动使用
func (m *Manger) AddRoom(r IRoom) {
	m.roomIdOnRoom.Store(r.GetId(), r)
}

// 加入房间，如果成功则返回 true
func (m *Manger) JoinRoom(player IPlayer, roomId int64) bool {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		return room.JoinRoom(player)
	}
	return false
}

// 退出房间，当房间内没有玩家时，移除房间
func (m *Manger) QuitRoom(player IPlayer, roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		room.QuitRoom(player)
	}
}

// 退出房间，当房间内没有玩家时，移除房间
func (m *Manger) QuitRoomByUserId(userId int64) {
	if value, ok := m.userOnRoom.Load(userId); ok {
		room := value.(IRoom)
		room.QuitRoom(room.GetPlayer(userId))
	}
}

// 移除房间并清理关联用户id
func (m *Manger) RemoveRoom(roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		for _, userId := range room.GetUserIds() {
			m.userOnRoom.Delete(userId)
		}
	}
	m.roomIdOnRoom.Delete(roomId)
}

// 获取房间列表
func (m *Manger) GetRooms() (list []IRoom) {
	m.roomIdOnRoom.Range(func(key, value any) bool {
		list = append(list, value.(IRoom))
		return true
	})
	return
}

// GetByRoomId 根据房间id获取房间， 对象,是否存在
func (m *Manger) GetByRoomId(roomId int64) IRoom {
	value, ok := m.roomIdOnRoom.Load(roomId)
	if ok {
		return value.(IRoom)
	}
	return nil
}

type State int

// 房间状态
const (
	Ready   State = iota // 准备状态, 未开始
	Running              // 运行状态, 已经开始
	Close                // 关闭状态, 已经结束
)

type IPlayer interface {
	UserId() int64
}

// DefaultPlayer 默认Player
type DefaultPlayer struct {
	Uid int64
}

func (d *DefaultPlayer) UserId() int64 {
	return d.Uid
}
