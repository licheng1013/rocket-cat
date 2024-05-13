package main

import (
	"log"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/room"
	"github.com/licheng1013/rocket-cat/router"
)

func main() {
	// 构建一个默认服务
	gateway := core.DefaultGateway()
	manager := room.NewManger()

	match := room.NewMatchQueue(2,func(players []room.IPlayer) {
		log.Println("匹配成功:", players)
		room := room.NewRoom(manager)
		for _, player := range players {
			room.JoinRoom(player)
		}
		room.Start(func() {
			//log.Println("帧同步")
			if room.GetLastSyncData().Len() == 0 {
				return
			}
			message := gateway.ToRouterData(1, 2, room.GetLastSyncData().GetList())
			gateway.Push(message)
		})
	})

	login := gateway.GetPlugin(core.LoginPluginId).(core.LoginInterface)
	// 添加一个路由
	gateway.Action(1, 1, func(ctx *router.Context) {
		var pos PosXY
			log.Println("收到:", pos)
			match.AddMatch(&room.DefaultPlayer{Uid: pos.UserId})
			ctx.Result(router.H{"userId": pos.UserId,"message":"等待其他玩家加入"})
			match.AddMatch(&room.DefaultPlayer{Uid: pos.UserId})
			ctx.Result(router.H{"userId": pos.UserId,"message":"等待其他玩家加入"})
		}
	})

	gateway.Action(1, 2, func(ctx *router.Context) {
		var pos PosXY
		_ = ctx.Message.Bind(&pos)
		//log.Println("收到:", pos)
		if r := manager.GetByUserId(pos.UserId); r != nil {
			r.(*room.SyncRoom).AddSyncData(pos)
		}
	})

	// 绑定路由
	gateway.Start(":10100")
}

type PosXY struct {
	X      int   `json:"x" form:"x"`
	Y      int   `json:"y" form:"y"`
	UserId int64 `json:"userId" form:"userId"`
}
