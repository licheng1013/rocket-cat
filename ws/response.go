package ws

func OK(ctx *Context, data any) {
	write(ctx, 0, "ok", data)
}

func Fail(ctx *Context, code int32, msg string) {
	write(ctx, code, msg, nil)
}

func write(ctx *Context, code int32, msg string, data any) {
	if ctx == nil || ctx.Session == nil {
		return
	}

	resp := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	if ctx.Packet != nil {
		resp.Cmd = ctx.Packet.Cmd
		resp.SubCmd = ctx.Packet.SubCmd
	}

	_ = ctx.Session.SendJSON(resp)
}
