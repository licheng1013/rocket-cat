package ws

import "fmt"

type Handler func(*Context)

type Middleware func(Handler) Handler

type Router struct {
	handlers    map[Route]Handler
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[Route]Handler),
	}
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

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
