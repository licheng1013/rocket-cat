package room

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoom(t *testing.T) {
	manger := NewManger()
	room := NewRoom(manger)
	assert.Equal(t, manger.JoinRoom(&player{Uid: 1}, room.GetId()), true)
	assert.Equal(t, manger.JoinRoom(&player{Uid: 1}, room.GetId()), true)
	assert.Equal(t, len(room.GetUserIds()), 1)
}
