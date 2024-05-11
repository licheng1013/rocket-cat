package room

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoomManager(t *testing.T) {
	user1 := &player{Uid: 1}
	user2 := &player{Uid: 2}
	// 创建房间
	manger := NewManger()
	room := NewRoom(manger)
	// 房间数
	assert.Equal(t, len(manger.GetRooms()), 1)
	// 玩家加入房间
	assert.Equal(t, manger.JoinRoom(user2, room.GetId()), true)
	assert.Equal(t, room.JoinRoom(user1), true)
	// 玩家数
	assert.Equal(t, len(room.GetUserIds()), 2)
	// 移除玩家
	room.QuitRoom(user1)
	assert.Equal(t, len(room.GetUserIds()), 1)
	// 移除房间
	manger.RemoveRoom(room.GetId())
	assert.Equal(t, len(manger.GetRooms()), 0)
}

type player struct {
	Uid int64
}

func (p *player) UserId() int64 {
	return p.Uid
}
