package chat

import "github.com/licheng1013/rocket-cat/ws"

func Register(router *ws.Router) {
	router.Register(Cmd, Query, QueryHandler)
	router.Register(Cmd, Create, CreateHandler)
	router.Register(Cmd, List, ListHandler)
}

func QueryHandler(ctx *ws.Context) {
	req, err := ws.Bind[QueryReq](ctx)
	if err != nil {
		return
	}

	resp, err := QueryService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}

func CreateHandler(ctx *ws.Context) {
	req, err := ws.Bind[CreateReq](ctx)
	if err != nil {
		return
	}

	var fromUid int64
	if ctx.Session != nil {
		fromUid = ctx.Session.Uid
	}
	resp, err := CreateService(req, fromUid)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}

func ListHandler(ctx *ws.Context) {
	req, err := ws.Bind[ListReq](ctx)
	if err != nil {
		return
	}

	resp, err := ListService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}
