package core

import (
	"fmt"
	"testing"
)

func TestLoginBody(t *testing.T) {
	t.Log("HelloWorld")
	l := &LoginBody{Action: LogoutByUserId, UserIds: []int64{1, 22}, State: true}
	data := l.ToMarshal()
	l2 := &LoginBody{}
	l2.ToUnmarshal(data)
	t.Log(l2)
}

func TestMap(t *testing.T) {
	login := &LoginPlugin{}
	login.Login(2, 1)
	fmt.Printf("login.ListUserId(): %v\n", login.ListUserId())
	fmt.Printf("login.ListSocketId(): %v\n", login.ListUserId())
}
