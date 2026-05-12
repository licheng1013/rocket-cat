package framesync

import "github.com/licheng1013/rocket-cat/ws"

const Cmd uint16 = 5001

const (
	Check     uint16 = 1
	JoinMatch uint16 = 2
	ExitMatch uint16 = 3
	Submit    uint16 = 4
	Push      uint16 = 6
)

var (
	CheckRoute     = ws.Merge(Cmd, Check)
	JoinMatchRoute = ws.Merge(Cmd, JoinMatch)
	ExitMatchRoute = ws.Merge(Cmd, ExitMatch)
	SubmitRoute    = ws.Merge(Cmd, Submit)
	PushRoute      = ws.Merge(Cmd, Push)
)
