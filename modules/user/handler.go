package user

import "github.com/licheng1013/rocket-cat/ws"

func Register(router *ws.Router) {
	router.Register(Cmd, Login, LoginHandler)
	router.Register(Cmd, Logout, LogoutHandler)
	router.Register(Cmd, Info, InfoHandler)
}

func LoginHandler(ctx *ws.Context) {
	req, err := ws.Bind[LoginReq](ctx)
	if err != nil {
		return
	}

	resp, err := LoginService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	if ctx.Session != nil {
		ctx.Session.Uid = resp.Uid
	}
	ws.OK(ctx, resp)
}

func LogoutHandler(ctx *ws.Context) {
	req, err := ws.Bind[LogoutReq](ctx)
	if err != nil {
		return
	}

	if err := LogoutService(req); err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	if ctx.Session != nil {
		ctx.Session.Uid = 0
	}
	ws.OK(ctx, map[string]bool{"ok": true})
}

func InfoHandler(ctx *ws.Context) {
	req, err := ws.Bind[InfoReq](ctx)
	if err != nil {
		return
	}

	resp, err := InfoService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}
