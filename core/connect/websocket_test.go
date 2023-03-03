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
	time.Sleep(1 * time.Second) // 需要等待1秒让客户端启动完成
	socket.SendMessage([]byte("广播消息-HelloWorld"))
	for i := 0; i < 2; i++ {
		select {
		case ok := <-channel:
			log.Println(ok)
		}
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
	go func() {
		for {
			_ = c.WriteMessage(websocket.BinaryMessage, []byte(message)) //和广播一起测试
			break
			//err := c.WriteMessage(websocket.BinaryMessage, []byte(message)) //此处用于多次数据发送
			//if err != nil {
			//	break
			//}
		}
	}()
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
