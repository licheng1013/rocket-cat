package main

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"runtime"
	"time"
)

func main() {

	gateway := core.DefaultGateway()
	go func() {
		// 打印线程数
		for true {
			fmt.Println("协程数 -> ", runtime.NumGoroutine())
			time.Sleep(3 * time.Second)
		}
	}()

	gateway.Router().AddAction(1, 1, func(ctx *router.Context) {
		var body core.LoginBody
		_ = ctx.Message.Bind(&body)
		log.Println("获取数据 -> ", string(ctx.Message.GetBody()))
		r := gateway.GetPlugin(core.LoginPluginId)
		login := r.(core.LoginInterface)
		if login.Login(body.UserId, ctx.SocketId) {
			fmt.Printf("login.ListUserId(): %v\n", login.ListUserId())
			login.SendAllUserMessage(ctx.Message.SetBody([]byte("用户")).GetBytesResult())
		}
		gateway.SendMessage(ctx.Message.SetBody([]byte("广播")).GetBytesResult())
		ctx.Message.SetBody([]byte("业务返回Hi->Ok->2"))
	})
	socket := &connect.WebSocket{}
	socket.Debug = false
	gateway.Start(connect.Addr, socket)
}
