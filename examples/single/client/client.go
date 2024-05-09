package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"log"
	"time"
)

// websocket连接塑胶
func main() {
	// 开始时间
	start := time.Now().UnixMilli()
	Single()
	time.Sleep(3 * time.Second)
	endTime := time.Now().UnixMilli()
	// 打印
	fmt.Printf("总耗时：%vms\n", endTime-start)
}

func Single() {
	maxSize := 10
	msgSize := 0
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+connect.Addr+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c := &core.LoginBody{UserId: 1}
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1), Body: c.ToMarshal()}
	jsonDecoder := decoder.JsonDecoder{}
	go func() {
		for {
			// 读取消息
			_, data, _ := conn.ReadMessage()
			msgSize++
			// 打印消息
			m := jsonDecoder.Decoder(data)
			log.Println("收到消息：", m.GetMerge(), string(m.GetBody()))
			if msgSize >= maxSize {
				return
			}
		}
	}()
	go func() {
		for i := 0; i < maxSize; i++ {
			// 发送消息
			_ = conn.WriteMessage(websocket.BinaryMessage, message.GetBytesResult())
		}
	}()
}
