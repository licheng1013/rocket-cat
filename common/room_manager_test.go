package common

import (
	"testing"
	"time"
)

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

	// 测试清理房间
	manger.RoomClear(1)

	time.Sleep(3 * time.Second)
	t.Log("sleep 3 second")
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
