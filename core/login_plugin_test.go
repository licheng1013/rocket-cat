package core

import (
    "github.com/licheng1013/rocket-cat/remote"
    "log"
    "testing"
)

func TestLoginBody(t *testing.T) {
	t.Log("HelloWorld")
    l := &LoginBody{LoginAction: LogoutByUserId,UserId: []int64{1,22},SocketId: []uint32{3,4}}
    data, _ := l.ToMarshal()

    l2 := &LoginBody{}
    _ = l2.ToUnmarshal(data)
	t.Log(l2)
}



func TestMap(t *testing.T) {

}