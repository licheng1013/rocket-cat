package room

import "github.com/licheng1013/rocket-cat/ws"

const Cmd uint16 = 3001

const (
	Query  uint16 = 1
	Create uint16 = 2
	List   uint16 = 5
)

var (
	QueryRoute  = ws.Merge(Cmd, Query)
	CreateRoute = ws.Merge(Cmd, Create)
	ListRoute   = ws.Merge(Cmd, List)
)
