package common

import (
	mrand "math/rand"
	"sync"
	"time"
)

// RoomManger 线程安全
type RoomManger struct {
	// 用户id - 房间id
	userOnRoom sync.Map
	// roomId - 房间具体房间
	roomIdOnRoom sync.Map
}

func NewRoomManger() *RoomManger {
	return &RoomManger{}
}

// GetUniqueRoomId 创建房间
func (m *RoomManger) GetUniqueRoomId() int64 {
	for {
		mrand.Seed(time.Now().UnixNano()) // 设置种子为当前时间戳
		n := mrand.Int63n(100000)
		if b := m.GetByRoomId(n); b == nil { // 直到不存在房间时赋予id
			return n
		}
	}
}

func (m *RoomManger) GetByUserId(userId int64) IRoom {
	value, ok := m.userOnRoom.Load(userId)
	if ok {
		return value.(IRoom)
	}
	return nil
}

// AddRoom 添加房间
func (m *RoomManger) AddRoom(r IRoom) {
	m.roomIdOnRoom.Store(r.GetRoomId(), r)
}

// JoinRoom 加入房间
func (m *RoomManger) JoinRoom(player IPlayer, roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		room.JoinRoom(player)
		m.userOnRoom.Store(player.UserId(), room)
	}
}

// QuitRoom 退出房间
func (m *RoomManger) QuitRoom(player IPlayer, roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		m.userOnRoom.Delete(player.UserId())
		room.QuitRoom(player)
	}
}

// RemoveRoom 移除房间并清理关联用户id
func (m *RoomManger) RemoveRoom(roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(IRoom)
		for _, userId := range room.GetUserIdList() {
			m.userOnRoom.Delete(userId)
		}
	}
	m.roomIdOnRoom.Delete(roomId)
}

// RoomClear 处理已经打开的房间并且30秒没有同步数据的房间 max 最大房间未动秒数
func (m *RoomManger) RoomClear(max int64) {
	go func() {
		for {
			m.roomIdOnRoom.Range(func(key, value any) bool {
				room := value.(IRoom)
				if room.GetRoomStatus() == Running {
					// 当超过max秒没有同步数据时，清理房间
					if room.HeartbeatTime()+max < time.Now().Unix() {
						m.RemoveRoom(room.GetRoomId())
					}
				}
				// 当房间状态为关闭时，清理房间
				if room.GetRoomStatus() == Close {
					m.RemoveRoom(room.GetRoomId())
				}
				return true
			})
			time.Sleep(time.Second)
		}
	}()
}

// ListRoom 获取房间列表
func (m *RoomManger) ListRoom() (list []IRoom) {
	m.roomIdOnRoom.Range(func(key, value any) bool {
		list = append(list, value.(IRoom))
		return true
	})
	return
}

// GetByRoomId 根据房间id获取房间， 对象,是否存在
func (m *RoomManger) GetByRoomId(roomId int64) IRoom {
	value, ok := m.roomIdOnRoom.Load(roomId)
	if ok {
		return value.(IRoom)
	}
	return nil
}

type RoomStatus int

// 房间状态
const (
	Ready   RoomStatus = iota // 准备状态, 未开始
	Running                   // 运行状态, 已经开始
	Close                     // 关闭状态, 已经结束
)

type IPlayer interface {
	UserId() int64
}
