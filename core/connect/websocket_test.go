package connect

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestWsServer(t *testing.T) {
	channel := make(chan int)
	socket := WebSocket{}
	go func() {
		socket.ListenBack(func(bytes []byte) []byte {
			return bytes
		})
		socket.ListenAddr(Addr)
	}()
	go WsClient(channel)
	time.Sleep(time.Second) //等待客户端完全启动
	socket.SendMessage([]byte("广播消息-HelloWorld"))
	select {
	case ok := <-channel:
		log.Println(ok)
	}
}

func WsClient(channel chan int) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	//go func() {
	//	for {
	//		// 2 标识字节消息
	//		err := c.WriteMessage(websocket.BinaryMessage, []byte(message))
	//		if err != nil {
	//			break
	//		}
	//	}
	//}()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("读取消息错误:", err)
			break
		}
		log.Println("获取数据:" + string(msg)) // 此处可能会打印多次,因为 channel <- 0 传输到通道也需要时间
		channel <- 0
	}
}
