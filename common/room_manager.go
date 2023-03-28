package common

import (
	mrand "math/rand"
	"sync"
	"time"
)

// RoomManger 房间管理器
var RoomManger = &roomManger{}

// 不允许外部直接创建
type roomManger struct {
	// 用户id - 房间
	userOnRoom sync.Map
	// roomId - 房间
	roomIdOnRoom sync.Map
}

// CreateRoom 创建房间
func (m *roomManger) CreateRoom() *Room {
	room := &Room{}
	for {
		mrand.Seed(time.Now().UnixNano()) // 设置种子为当前时间戳
		n := mrand.Int63n(100000)
		if b := m.GetByRoomId(n); b == nil { // 直到不存在房间时赋予id
			room.RoomId = n
			break
		}
	}
	return room
}

func (m *roomManger) GetByUserId(userId int64) *Room {
	value, ok := m.userOnRoom.Load(userId)
	if ok {
		return value.(*Room)
	}
	return nil
}

// AddRoom 添加房间
func (m *roomManger) AddRoom(r *Room) {
	m.roomIdOnRoom.Store(r.RoomId, r)
}

// PlayerJoinRoom 加入房间
func (m *roomManger) PlayerJoinRoom(player Player, roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(*Room)
		room.UserList = append(room.UserList, player)
		m.userOnRoom.Store(player.UserId(), room)
	}
}

// PlayerQuitRoom 退出房间
func (m *roomManger) PlayerQuitRoom(player Player, roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(*Room)
		m.userOnRoom.Delete(player.UserId)
		var delIndex int64
		for i, item := range room.UserList {
			if item.UserId() == player.UserId() {
				delIndex = int64(i)
				break
			}
		}
		room.UserList = append(room.UserList[:delIndex], room.UserList[delIndex+1:]...)
	}
}

// RemoveRoom 移除房间并清理关联用户id
func (m *roomManger) RemoveRoom(roomId int64) {
	if value, ok := m.roomIdOnRoom.Load(roomId); ok {
		room := value.(*Room)
		for _, user := range room.UserList {
			m.userOnRoom.Delete(user.UserId())
		}
	}
	m.roomIdOnRoom.Delete(roomId)
}

// RoomClear 处理已经打开的房间并且30秒没有同步数据的房间 max 最大房间未动秒数
func (m *roomManger) RoomClear(max int64) {
	m.roomIdOnRoom.Range(func(key, value any) bool {
		room := value.(*Room)
		if room.RoomStatus == Open {
			if room.LastSyncTime+max < time.Now().Unix() {
				room.RoomStatus = Close
			}
		}
		if room.RoomStatus == Close {
			m.RemoveRoom(room.RoomId)
		}
		return true
	})
}

// ListRoom 获取房间列表
func (m *roomManger) ListRoom() (list []Room) {
	m.roomIdOnRoom.Range(func(key, value any) bool {
		list = append(list, value.(Room))
		return true
	})
	return
}

// GetByRoomId 根据房间id获取房间， 对象,是否存在
func (m *roomManger) GetByRoomId(roomId int64) *Room {
	value, ok := m.roomIdOnRoom.Load(roomId)
	if ok {
		return value.(*Room)
	}
	return nil
}

type RoomStatus int

const (
	Open  RoomStatus = iota // 房间开启
	Close                   // 房间关闭 -> 需要被清理线程进行清理删除掉了
)

// Room 房间
type Room struct {
	// Id
	RoomId int64
	// 用户Ids
	UserList []Player
	// 房间状态
	RoomStatus
	// 同步数据
	List []*SafeMap
	// 十位时间戳
	LastSyncTime int64
}

// Start 进行房间的帧同步，以每秒60帧为例，每1/60秒执行一次
func (r *Room) Start(f func()) {
	r.StartCustom(f,time.Second/60)
}

// StartCustom 以每秒60帧为例，delay = time.Second/60 为每一帧的执行时间 = 1/60m秒
func (r *Room) StartCustom(f func(),delay time.Duration) {
	// 使用 common.SyncManager 进行帧同步
	// 帧同步数据
	r.RoomStatus = Open
	manager := NewFrameSyncManager(60,delay)
	manager.Start()
	go func() {
		for true {
			if r == nil || r.RoomStatus == Close {
				return
			}
			// 执行每一帧
			manager.WaitNextFrame(f)
			r.List = append(r.List, NewSafeMap())
		}
	}()
}


// AddSyncData 添加同步数据
func (r *Room) AddSyncData(userId int64, value any) {
	r.LastSyncTime = time.Now().Unix()
	if len(r.List) == 0 {
		return
	}
	safeMap := r.List[len(r.List)-1]
	if safeMap != nil {
		safeMap.Set(userId, value)
	}
}

// GetLastSyncData 获取最后帧的同步数据
func (r *Room) GetLastSyncData() *SafeMap {
	if len(r.List) == 0 {
		return NewSafeMap()
	}
	return r.List[len(r.List)-1]
}

// UserIds 获取所有用户Id
func (r *Room) UserIds() (list []int64) {
	for _, player := range r.UserList {
		list = append(list, player.UserId())
	}
	return
}

type Player interface {
	UserId() int64
}
