package chat

import "github.com/licheng1013/rocket-cat/ws"

// Register 显式注册聊天模块路由。
func Register(router *ws.Router) {
	router.Register(Cmd, Query, QueryHandler)
	router.Register(Cmd, Create, CreateHandler)
	router.Register(Cmd, List, ListHandler)
}

// QueryHandler 处理聊天消息查询请求。
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

// CreateHandler 处理聊天消息创建请求。
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

// ListHandler 处理聊天消息列表请求。
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
