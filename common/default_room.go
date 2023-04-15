package common

import (
	"sync"
	"time"
)

// IRoom 房间接口, 其默认实现为DefaultRoom, 所以继承DefaultRoom使用即可
type IRoom interface {
	// GetRoomId 获取房间id
	GetRoomId() int64
	// GetMaxUser 获取房间最大人数
	GetMaxUser() int64
	// GetRoomStatus 获取房间状态
	GetRoomStatus() RoomStatus
	// GetUserIdList 获取房间内所有玩家Id
	GetUserIdList() []int64
	// GetPlayerList 获取房间内所有玩家
	GetPlayerList() []IPlayer
	// JoinRoom 加入房间，请通过RoomManager使用
	JoinRoom(player IPlayer)
	// QuitRoom 退出房间，请通过RoomManager使用
	QuitRoom(player IPlayer)
	// HeartbeatTime 房间上次心跳时间,如果在一定时间内没有心跳则清理房间,10位时间戳
	HeartbeatTime() int64
	// GetPlayer 获取某个玩家
	GetPlayer(userId int64) IPlayer
}

// DefaultRoom 默认房间实现,请继承此结构体,并重写方法,线程安全
type DefaultRoom struct {
	// 房间id
	RoomId int64
	// 用户Ids
	userList []IPlayer
	// 房间状态
	RoomStatus
	// 创建时间十位时间戳
	CreateTime int64
	// 最大人数
	MaxUser int64
	// 心跳时间
	Heartbeat int64
	// 锁
	sync.Mutex
}

func (d *DefaultRoom) GetPlayer(userId int64) IPlayer {
	defer d.Unlock()
	d.Lock()
	for _, player := range d.userList {
		if player.UserId() == userId {
			return player
		}
	}
	return nil
}

func (d *DefaultRoom) GetRoomId() int64 {
	return d.RoomId
}

func (d *DefaultRoom) GetMaxUser() int64 {
	return d.MaxUser
}

func (d *DefaultRoom) GetRoomStatus() RoomStatus {
	return d.RoomStatus
}

func (d *DefaultRoom) GetUserIdList() (list []int64) {
	defer d.Unlock()
	d.Lock()
	for _, player := range d.userList {
		list = append(list, player.UserId())
	}
	return
}
func (d *DefaultRoom) GetPlayerList() []IPlayer {
	defer d.Unlock()
	d.Lock()
	return d.userList
}
func (d *DefaultRoom) JoinRoom(player IPlayer) {
	defer d.Unlock()
	d.Lock()
	d.userList = append(d.userList, player)
}
func (d *DefaultRoom) QuitRoom(player IPlayer) {
	defer d.Unlock()
	d.Lock()
	var delIndex int64
	for i, item := range d.userList {
		if item.UserId() == player.UserId() {
			delIndex = int64(i)
			break
		}
	}
	d.userList = append(d.userList[:delIndex], d.userList[delIndex+1:]...)
}

func (d *DefaultRoom) HeartbeatTime() int64 {
	return d.Heartbeat
}

// SyncRoom 帧同步房间
type SyncRoom struct {
	DefaultRoom
	// 同步数据,索引为帧号
	List []*SafeList
	// 创建时间十位时间戳
	CreateTime int64
}

func NewRoom(roomId int64) *SyncRoom {
	r := &SyncRoom{CreateTime: time.Now().Unix()}
	r.RoomId = roomId
	r.RoomStatus = Ready
	return r
}

// Start 进行房间的帧同步，以每秒60帧为例，每1/60秒执行一次
func (r *SyncRoom) Start(f func()) {
	r.StartCustom(f, time.Second/60)
}

// StartCustom 以每秒60帧为例，delay = time.Second/60 为每一帧的执行时间 = 1/60m秒
func (r *SyncRoom) StartCustom(f func(), delay time.Duration) {
	// 使用 common.SyncManager 进行帧同步
	// 帧同步数据
	r.RoomStatus = Running
	manager := NewFrameSyncManager(60, delay)
	manager.Start()
	go func() {
		for true {
			if r == nil || r.RoomStatus == Close {
				return
			}
			// 执行每一帧
			manager.WaitNextFrame(f)
			r.List = append(r.List, &SafeList{})
		}
	}()
}

// AddSyncData 添加同步数据
func (r *SyncRoom) AddSyncData(value any) {
	r.Heartbeat = time.Now().Unix()
	if len(r.List) == 0 {
		return
	}
	safeMap := r.List[len(r.List)-1]
	if safeMap != nil {
		safeMap.Add(value)
	}
}

// GetLastSyncData 获取最后帧的同步数据
func (r *SyncRoom) GetLastSyncData() *SafeList {
	if len(r.List) == 0 {
		return &SafeList{}
	}
	return r.List[len(r.List)-1]
}

// GetUserIdList  获取所有用户Id
func (r *SyncRoom) GetUserIdList() (list []int64) {
	for _, player := range r.userList {
		list = append(list, player.UserId())
	}
	return
}
