package room

import (
	"sync"
	"time"
)

// RoomManger 线程安全
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

// AddRoom 添加房间
func (m *Manger) AddRoom(r IRoom) {
	m.roomIdOnRoom.Store(r.GetRoomId(), r)
}

// 加入房间，如果成功则返回 true
func (m *Manger) JoinRoom(player IPlayer, roomId int64) bool {
	var join bool
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		room.JoinRoom(player)
		m.userOnRoom.Store(player.UserId(), room)
		join = true
	}
	return join
}

// 退出房间
func (m *Manger) QuitRoom(player IPlayer, roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		m.userOnRoom.Delete(player.UserId())
		room.QuitRoom(player)
	}
}

// RemoveRoom 移除房间并清理关联用户id
func (m *Manger) RemoveRoom(roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		for _, userId := range room.GetUserIdList() {
			m.userOnRoom.Delete(userId)
		}
		// 清理房间
		room.ClearRoom()
	}
	m.roomIdOnRoom.Delete(roomId)
}

// RoomClear 处理已经打开的房间并且30秒没有同步数据的房间 max 最大房间未动秒数
func (m *Manger) RoomClear(max int64) {
	go func() {
		for {
			m.roomIdOnRoom.Range(func(key, value any) bool {
				room := value.(IRoom)
				// 当超过max秒没有同步数据时，清理房间
				if room.HeartbeatTime()+max < time.Now().Unix() {
					m.RemoveRoom(room.GetRoomId())
				}
				return true
			})
			time.Sleep(time.Second)
		}
	}()
}

// ListRoom 获取房间列表
func (m *Manger) ListRoom() (list []IRoom) {
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

type Status int

// 房间状态
const (
	Ready   Status = iota // 准备状态, 未开始
	Running               // 运行状态, 已经开始
	Close                 // 关闭状态, 已经结束
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
