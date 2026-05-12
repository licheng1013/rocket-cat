package ws

import "fmt"

// Handler 表示统一的 WS 业务处理函数。
type Handler func(*Context)

// Middleware 表示对 Handler 的包装函数。
type Middleware func(Handler) Handler

// Router 保存路由到处理函数的映射。
type Router struct {
	handlers    map[Route]Handler // 路由处理函数表
	middlewares []Middleware      // 全局中间件列表
}

// NewRouter 创建一个空路由器。
func NewRouter() *Router {
	return &Router{
		handlers: make(map[Route]Handler),
	}
}

// Use 添加全局中间件。
func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

// Register 显式注册主命令和子命令对应的处理函数。
func (r *Router) Register(cmd uint16, sub uint16, handler Handler) {
	if handler == nil {
		panic("ws: nil handler")
	}

	route := Merge(cmd, sub)
	if _, ok := r.handlers[route]; ok {
		panic(fmt.Sprintf("ws: duplicate route cmd=%d subCmd=%d", cmd, sub))
	}

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	r.handlers[route] = handler
}

// Dispatch 根据请求包中的命令分发到对应处理函数。
func (r *Router) Dispatch(ctx *Context) {
	if ctx == nil || ctx.Packet == nil {
		Fail(ctx, 400, "empty packet")
		return
	}

	handler, ok := r.handlers[Merge(ctx.Packet.Cmd, ctx.Packet.SubCmd)]
	if !ok {
		Fail(ctx, 404, "route not found")
		return
	}

	handler(ctx)
}
