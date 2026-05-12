package framesync

import (
	"sync"

	"github.com/licheng1013/rocket-cat/ws"
)

var sessions sync.Map

// Register explicitly registers frame sync routes.
func Register(router *ws.Router) {
	router.Register(Cmd, Check, CheckHandler)
	router.Register(Cmd, JoinMatch, JoinMatchHandler)
	router.Register(Cmd, ExitMatch, ExitMatchHandler)
	router.Register(Cmd, Submit, SubmitHandler)
}

// CheckHandler handles sync room existence checks.
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

// JoinMatchHandler handles joining the matching queue.
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

// ExitMatchHandler handles exiting the matching queue.
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

// SubmitHandler handles client input submission.
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

func bindSession(clientId string, session *ws.Session) {
	if clientId == "" || session == nil {
		return
	}
	sessions.Store(clientId, session)
}

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
