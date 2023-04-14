package common

import (
	"testing"
)

func TestRoom(t *testing.T) {
	manger := NewRoomManger()
	roomId := manger.GetUniqueRoomId()
	room := NewRoom(roomId)
	manger.AddRoom(room)
	manger.JoinRoom(&Player{Uid: 1}, room.GetRoomId())
	manger.JoinRoom(&Player{Uid: 2}, room.GetRoomId())
	t.Log(room.GetUserIdList())
}
