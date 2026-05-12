package user

import "github.com/licheng1013/rocket-cat/ws"

const Cmd uint16 = 1001

const (
	Login  uint16 = 1
	Logout uint16 = 2
	Info   uint16 = 3
)

var (
	LoginRoute  = ws.Merge(Cmd, Login)
	LogoutRoute = ws.Merge(Cmd, Logout)
	InfoRoute   = ws.Merge(Cmd, Info)
)
