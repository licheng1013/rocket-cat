package room

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoom(t *testing.T) {
	manger := NewManger()
	roomId := manger.GetUniqueRoomId()
	room := NewRoom(roomId)
	manger.AddRoom(room)
	assert.Equal(t, manger.JoinRoom(&player{Uid: 1}, room.GetRoomId()), true)
	assert.Equal(t, manger.JoinRoom(&player{Uid: 1}, room.GetRoomId()), true)
	assert.Equal(t, len(room.GetUserIdList()), 1)
}
