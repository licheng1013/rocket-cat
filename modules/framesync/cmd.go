package framesync

import "github.com/licheng1013/rocket-cat/ws"

// Cmd 是帧同步模块主命令。
const Cmd uint16 = 5001

const (
	Check     uint16 = 1 // 检查已有同步房间
	JoinMatch uint16 = 2 // 加入匹配
	ExitMatch uint16 = 3 // 退出匹配
	Submit    uint16 = 4 // 提交输入快照
	Push      uint16 = 6 // 服务端推送
)

var (
	CheckRoute     = ws.Merge(Cmd, Check)     // 检查已有同步房间路由
	JoinMatchRoute = ws.Merge(Cmd, JoinMatch) // 加入匹配路由
	ExitMatchRoute = ws.Merge(Cmd, ExitMatch) // 退出匹配路由
	SubmitRoute    = ws.Merge(Cmd, Submit)    // 提交输入快照路由
	PushRoute      = ws.Merge(Cmd, Push)      // 服务端推送路由
)
