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
	// 开始时间
	start := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		go Single()
	}
	time.Sleep(5 * time.Second)
	endTime := time.Now().UnixMilli()
	// 打印
	fmt.Printf("总耗时：%vms\n", endTime-start)
}

func Single() {
	max := 10000
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+connect.Addr+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c := &core.LoginBody{UserId: 1}
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 2), Body: c.ToMarshal()}
	var count int
	go func() {
		// 开始时间
		for true {
			// 读取消息
			_, _, err := conn.ReadMessage()
			count++
			if count >= max {
				return
			}
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	start := time.Now().UnixMilli()
	for true {
		if count >= max {
			fmt.Printf("收到消息：%v,%v\n", count, time.Now().UnixMilli()-start)
			return
		}
		// 发送消息
		err = conn.WriteMessage(websocket.BinaryMessage, message.GetBytesResult())
		if err != nil {
			log.Println(err)
			return
		}
	}
}
