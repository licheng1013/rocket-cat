package main

import (
	"log"
	"time"

	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/room"
	"github.com/licheng1013/rocket-cat/router"
)

// 构建一个默认服务
var gateway = core.DefaultGateway()
var manager = room.NewManger()
var login = gateway.GetPlugin(core.LoginPluginId).(core.LoginInterface)

func main() {

	go func() {
		for {
			time.Sleep(time.Second * 5)
			list := manager.GetRooms()
			for _, item := range list {
				unix := time.Now().Unix()
				// 如果超过10秒没有更新时间，就关闭房间
				if unix-item.GetUpdateTime() > 10 && item.GetState() == room.Running {
					log.Println("关闭房间:", item.GetId())
					manager.RemoveRoom(item.GetId())
				}
			}
		}
	}()

	match := room.NewMatchQueue(2, func(players []room.IPlayer) {
		log.Println("匹配成功:", players)
		newRoom := room.NewRoom(manager)
		for _, player := range players {
			newRoom.JoinRoom(player)
		}
		newRoom.State = room.Running
		newRoom.Start(func() {
			//log.Println("帧同步")
			if newRoom.GetLastSyncData().Len() == 0 {
				return
			}
			message := gateway.ToRouterData(1, 2, newRoom.GetLastSyncData().GetList())
			gateway.Push(message)
		})
	})

	gateway.AddCloseHook(func(socketId uint32) {
		if userId := login.(*core.LoginPlugin).GetUserIdBySocketId(socketId); userId != 0 {
			manager.QuitRoomByUserId(userId)
		}
	})

	// 添加一个路由
	gateway.Action(1, 1, func(ctx *router.Context) {
		var pos PosXY
		_ = ctx.Message.Bind(&pos)
		if login.Login(pos.UserId, ctx.SocketId) {
			log.Println("收到:", pos)
			match.AddMatch(&room.DefaultPlayer{Uid: pos.UserId})
			ctx.Result(router.H{"userId": pos.UserId, "message": "等待其他玩家加入"})
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
