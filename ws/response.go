package ws

// OK 向客户端返回成功响应。
func OK(ctx *Context, data any) {
	write(ctx, 0, "ok", data)
}

// Fail 向客户端返回失败响应。
func Fail(ctx *Context, code int32, msg string) {
	write(ctx, code, msg, nil)
}

// write 写入统一响应包。
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
