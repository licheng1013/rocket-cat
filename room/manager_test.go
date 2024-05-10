package room

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoomManager(t *testing.T) {
	// 创建房间
	manger := NewManger()
	roomId := manger.GetUniqueRoomId()
	room := NewRoom(roomId)
	assert.Equal(t, room.GetId(), roomId)
	assert.Equal(t, manger.JoinRoom(&player{}, roomId), false)

	manger.AddRoom(room)
	assert.Equal(t, manger.JoinRoom(&player{}, roomId), true)

	assert.NotEqual(t, manger.GetByRoomId(roomId), nil)

	assert.Equal(t, len(manger.GetRooms()), 1)

	assert.Equal(t, len(room.GetUserIds()), 1)
	manger.QuitRoom(&player{}, roomId)
	assert.Equal(t, len(room.GetUserIds()), 0)

	manger.RemoveRoom(roomId)
	assert.Equal(t, len(manger.GetRooms()), 0)
}

type player struct {
	Uid int64
}

func (p *player) UserId() int64 {
	return p.Uid
}
