package room

import (
	"github.com/licheng1013/rocket-cat/common"
	"time"
)

// SyncRoom 帧同步房间
type SyncRoom struct {
	DefaultRoom
	// 同步数据,索引为帧号
	List []*common.SafeList
}

func NewRoom(manager *Manger) *SyncRoom {
	r := &SyncRoom{}
	r.CreateTime = time.Now().Unix()
	r.RoomId = manager.GetUniqueRoomId()
	r.RoomId = manager.GetUniqueRoomId()
	r.Status = Ready
	r.manager = manager
	manager.AddRoom(r)
	r.manager = manager
	manager.AddRoom(r)
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
	r.Status = Running
	manager := NewFrameSyncManager(60, delay)
	manager.Start()
	go func() {
		for {
			if r == nil || r.Status == Close {
				return
			}
			// 执行每一帧
			manager.WaitNextFrame(f)
			r.List = append(r.List, &common.SafeList{})
		}
	}()
}

// AddSyncData 添加同步数据
func (r *SyncRoom) AddSyncData(value any) {
	if len(r.List) == 0 {
		return
	}
	safeMap := r.List[len(r.List)-1]
	if safeMap != nil {
		safeMap.Add(value)
	}
}

// GetLastSyncData 获取最后帧的同步数据
func (r *SyncRoom) GetLastSyncData() *common.SafeList {
	if len(r.List) == 0 {
		return &common.SafeList{}
	}
	return r.List[len(r.List)-1]
}
