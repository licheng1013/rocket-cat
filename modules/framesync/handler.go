package framesync

import (
	"sync"

	"github.com/licheng1013/rocket-cat/ws"
)

// sessions 保存客户端 ID 到当前连接会话的映射。
var sessions sync.Map

// Register 显式注册帧同步模块路由。
func Register(router *ws.Router) {
	router.Register(Cmd, Check, CheckHandler)
	router.Register(Cmd, JoinMatch, JoinMatchHandler)
	router.Register(Cmd, ExitMatch, ExitMatchHandler)
	router.Register(Cmd, Submit, SubmitHandler)
}

// CheckHandler 处理已有同步房间检查请求。
func CheckHandler(ctx *ws.Context) {
	req, err := ws.Bind[CheckReq](ctx)
	if err != nil {
		return
	}

	bindSession(req.ClientId, ctx.Session)
	resp, err := CheckService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}

// JoinMatchHandler 处理加入匹配请求。
func JoinMatchHandler(ctx *ws.Context) {
	req, err := ws.Bind[JoinMatchReq](ctx)
	if err != nil {
		return
	}

	bindSession(req.ClientId, ctx.Session)
	resp, event, err := JoinMatchService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
	if event != nil {
		pushToPlayers(event.Players, event)
	}
}

// ExitMatchHandler 处理退出匹配请求。
func ExitMatchHandler(ctx *ws.Context) {
	req, err := ws.Bind[ExitMatchReq](ctx)
	if err != nil {
		return
	}

	resp, err := ExitMatchService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}

// SubmitHandler 处理客户端输入提交请求。
func SubmitHandler(ctx *ws.Context) {
	req, err := ws.Bind[SubmitReq](ctx)
	if err != nil {
		return
	}

	bindSession(req.ClientId, ctx.Session)
	resp, event, err := SubmitService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
	if event != nil {
		pushToPlayers(event.Players, event)
	}
}

// bindSession 绑定客户端 ID 与当前连接会话。
func bindSession(clientId string, session *ws.Session) {
	if clientId == "" || session == nil {
		return
	}
	sessions.Store(clientId, session)
}

// pushToPlayers 将事件推送给指定客户端列表。
func pushToPlayers(players []string, event *PushEvent) {
	for _, player := range players {
		value, ok := sessions.Load(player)
		if !ok {
			continue
		}
		session, ok := value.(*ws.Session)
		if !ok {
			continue
		}
		_ = session.SendJSON(ws.Response{
			Cmd:    Cmd,
			SubCmd: Push,
			Code:   0,
			Msg:    "push",
			Data:   event,
		})
	}
}
