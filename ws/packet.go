package ws

import "encoding/json"

// Route 表示由主命令和子命令合并后的路由值。
type Route uint32

// Merge 将主命令和子命令合并为路由值。
func Merge(cmd, sub uint16) Route {
	return Route(uint32(cmd)<<16 | uint32(sub))
}

// Cmd 从路由值中解析主命令。
func Cmd(route Route) uint16 {
	return uint16(route >> 16)
}

// Sub 从路由值中解析子命令。
func Sub(route Route) uint16 {
	return uint16(route)
}

// Packet 表示客户端发送的统一协议包。
type Packet struct {
	Cmd    uint16          `json:"cmd"`    // 主命令
	SubCmd uint16          `json:"subCmd"` // 子命令
	Data   json.RawMessage `json:"data"`   // 业务数据
}

// Response 表示服务端返回的统一响应包。
type Response struct {
	Cmd    uint16 `json:"cmd"`            // 主命令
	SubCmd uint16 `json:"subCmd"`         // 子命令
	Code   int32  `json:"code"`           // 状态码
	Msg    string `json:"msg"`            // 状态消息
	Data   any    `json:"data,omitempty"` // 响应数据
}
