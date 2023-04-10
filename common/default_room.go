package common

import "time"

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
	// LastSyncTime 房间上次心跳时间,如果在一定时间内没有心跳则清理房间
	LastSyncTime() int64
}

// Room 房间
type Room struct {
	// Id
	RoomId int64
	// 用户Ids
	UserList []IPlayer
	// 房间状态
	RoomStatus
	// 同步数据,索引为帧号
	List []*SafeList
	// 十位时间戳
	lastSyncTime int64
	// 创建时间十位时间戳
	CreateTime int64
}

func NewRoom(roomId int64) *Room {
	return &Room{RoomId: roomId, RoomStatus: Ready, CreateTime: time.Now().Unix()}
}

func (r *Room) LastSyncTime() int64 {
	return r.lastSyncTime
}

func (r *Room) QuitRoom(player IPlayer) {
	var delIndex int64
	for i, item := range r.UserList {
		if item.UserId() == player.UserId() {
			delIndex = int64(i)
			break
		}
	}
	r.UserList = append(r.UserList[:delIndex], r.UserList[delIndex+1:]...)
}

func (r *Room) JoinRoom(player IPlayer) {
	r.UserList = append(r.UserList, player)
}

func (r *Room) GetPlayerList() []IPlayer {
	return r.UserList
}

func (r *Room) GetRoomId() int64 {
	return r.RoomId
}

func (r *Room) GetMaxUser() int64 {
	return 3
}

func (r *Room) GetRoomStatus() RoomStatus {
	return Ready
}

// Start 进行房间的帧同步，以每秒60帧为例，每1/60秒执行一次
func (r *Room) Start(f func()) {
	r.StartCustom(f, time.Second/60)
}

// StartCustom 以每秒60帧为例，delay = time.Second/60 为每一帧的执行时间 = 1/60m秒
func (r *Room) StartCustom(f func(), delay time.Duration) {
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
func (r *Room) AddSyncData(value any) {
	r.lastSyncTime = time.Now().Unix()
	if len(r.List) == 0 {
		return
	}
	safeMap := r.List[len(r.List)-1]
	if safeMap != nil {
		safeMap.Add(value)
	}
}

// GetLastSyncData 获取最后帧的同步数据
func (r *Room) GetLastSyncData() *SafeList {
	if len(r.List) == 0 {
		return &SafeList{}
	}
	return r.List[len(r.List)-1]
}

// GetUserIdList  获取所有用户Id
func (r *Room) GetUserIdList() (list []int64) {
	for _, player := range r.UserList {
		list = append(list, player.UserId())
	}
	return
}
