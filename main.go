package main

import (
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/router"
)

func main() {
	gateway := core.DefaultGateway()
	gateway.Action(1, 1, func(ctx *router.Context) {
		ctx.Message.SetBody([]byte("业务返回Hi->Ok->2"))
	})
	gateway.Start(connect.Addr)
}
