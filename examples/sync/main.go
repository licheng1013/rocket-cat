package main

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"time"
)

func main() {
	// 构建一个默认服务
	gateway := core.DefaultGateway()
	var group = gateway.Router()
	login := gateway.GetPlugin(core.LoginPluginId).(core.LoginInterface)
	// 添加一个路由
	group.Action(1, 1, func(ctx *router.Context) {
		var pos PosXY
		_ = ctx.Message.Bind(&pos)
		if login.Login(pos.UserId, ctx.SocketId) {
			log.Println("收到:", pos)
			ctx.Result(map[string]interface{}{"userId": pos.UserId})
		}
	})

	data := common.SafeList{}
	group.Action(1, 2, func(ctx *router.Context) {
		var pos PosXY
		_ = ctx.Message.Bind(&pos)
		//log.Println("收到:", pos)
		data.Add(pos)
	})
	jsonDecoder := decoder.JsonDecoder{}
	gateway.SetDecoder(jsonDecoder)

	go func() {
		// 每 16 毫秒发送一次消息
		for {
			time.Sleep(time.Second / 60)
			if data.Len() == 0 {
				continue
			}
			message := jsonDecoder.Data(1, 2, data.GetList())
			data = common.SafeList{}
			gateway.Push(jsonDecoder.Encode(message))
		}
	}()

	// 绑定路由
	gateway.Start(":10100")
}

type PosXY struct {
	X      int   `json:"x" form:"x"`
	Y      int   `json:"y" form:"y"`
	UserId int64 `json:"userId" form:"userId"`
}
