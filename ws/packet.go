package ws

import "encoding/json"

type Route uint32

func Merge(cmd, sub uint16) Route {
	return Route(uint32(cmd)<<16 | uint32(sub))
}

func Cmd(route Route) uint16 {
	return uint16(route >> 16)
}

func Sub(route Route) uint16 {
	return uint16(route)
}

type Packet struct {
	Cmd    uint16          `json:"cmd"`
	SubCmd uint16          `json:"subCmd"`
	Data   json.RawMessage `json:"data"`
}

type Response struct {
	Cmd    uint16 `json:"cmd"`
	SubCmd uint16 `json:"subCmd"`
	Code   int32  `json:"code"`
	Msg    string `json:"msg"`
	Data   any    `json:"data,omitempty"`
}
