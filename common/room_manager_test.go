package common

import "testing"

func TestRoomManager(t *testing.T) {
	// 创建房间
	manger := NewRoomManger()
	id := manger.GetUniqueRoomId()

	manger.AddRoom(NewRoom(id))
	manger.AddRoom(NewRoom(manger.GetUniqueRoomId()))

	manger.PlayerJoinRoom(Player{}, id)

	// 列出所有房间
	for _, room := range manger.ListRoom() {
		t.Log(room.GetRoomId())
		t.Log(room.GetUserIdList())
	}

	manger.RemoveRoom(id)
	// 列出所有房间
	for _, room := range manger.ListRoom() {
		t.Log(room.GetRoomId())
	}

}

type Player struct {
}

func (p Player) UserId() int64 {
	return 1000
}
