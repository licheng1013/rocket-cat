package common

import (
	"testing"
)

func TestRoomManager(t *testing.T) {
	// 创建房间
	manger := NewRoomManger()
	roomId := manger.GetUniqueRoomId()
	room := NewRoom(roomId)
	manger.AddRoom(room)
	manger.JoinRoom(&Player{}, roomId)
	t.Log(room.GetUserIdList())
	t.Log(manger.ListRoom())
	manger.QuitRoom(&Player{}, roomId)
	// 删除房间
	manger.RemoveRoom(roomId)
	// 打印
	t.Log(manger.ListRoom())
}

type Player struct {
}

func (p Player) UserId() int64 {
	return 1000
}
