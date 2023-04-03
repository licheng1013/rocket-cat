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
	Single()
	endTime := time.Now().UnixMilli()
	// 打印
	fmt.Printf("总耗时：%vms\n", endTime-start)
}

func Single() {
	max := 1000000
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+connect.Addr+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c := &core.LoginBody{UserId: 1}
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1), Body: c.ToMarshal()}
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

func ManyTest() {
	countChan := make(chan int, 100000)
	threadCount := 10
	for i := 0; i < threadCount; i++ {
		body1 := core.LoginBody{UserId: int64(i + 1)}
		go ConnServer(body1, countChan)
	}
	i := threadCount * 1000
	var c int64
	// 开始时间
	start := time.Now().UnixMilli()
	for range countChan {
		c++
		log.Println("收到消息：", c)
		if c >= int64(i) {
			break
		}
	}
	// 结束时间并打印
	end := time.Now().UnixMilli()
	fmt.Println("耗时：", end-start, "ms")
}

func ConnServer(body core.LoginBody, count chan int) {
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+connect.Addr+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1), Body: body.ToMarshal()}
	go func() {
		for true {
			// 读取消息
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			// 打印消息
			//fmt.Printf("收到消息: %s - %v\n", p, body.UserId)
			count <- 1
		}
	}()
	for true {
		// 发送消息
		err = conn.WriteMessage(websocket.BinaryMessage, message.GetBytesResult())
		if err != nil {
			log.Println(err)
			return
		}
	}
}
