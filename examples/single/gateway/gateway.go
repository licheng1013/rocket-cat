package main

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/router"
)

func main() {
	gateway := core.DefaultGateway()
	gateway.Router().Action(1, 1, func(ctx *router.Context) {
		var body core.LoginBody
		_ = ctx.Message.Bind(&body)
		common.CatLog.Println("获取数据:", body)
		r := gateway.GetPlugin(core.LoginPluginId)
		login := r.(core.LoginInterface)
		if login.Login(body.UserId, ctx.SocketId) {
			fmt.Println("登入用户:", login.GetUserIds())
			//login.Push(gateway.ToRouterData(1, 1, []byte("HelloWorld")))
		}
		//gateway.Push(gateway.ToRouterData(1, 1, []byte("HelloWorld")))
		ctx.Result([]byte("业务返回Hi->Ok->2"))
	})
	socket := &connect.WebSocket{}
	socket.Debug = false
	gateway.SetSocket(socket)
	gateway.Start(connect.Addr)
}
