package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/messages"
	"log"
	"time"
)

// 测试多用户连接
func main() {

	body1 := core.LoginBody{UserId: 1}
	body2 := core.LoginBody{UserId: 2}
	go ConnServer(body1)
	time.Sleep(time.Second)
	ConnServer(body2)
}

func ConnServer(body core.LoginBody) {
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+connect.Addr+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1), Body: body.ToMarshal()}
	go func() {
		for true {
			// 读取消息
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			// 打印消息
			fmt.Printf("收到消息: %s - %v\n", p, body.UserId)
		}
	}()
	i := 0
	for true {
		if i == 0 {
			// 发送消息
			err = conn.WriteMessage(websocket.BinaryMessage, message.GetBytesResult())
			if err != nil {
				log.Println(err)
				return
			}
		}
		i++
	}
}
