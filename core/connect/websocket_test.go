package connect

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
)

func TestWsServer(t *testing.T) {
	channel := make(chan int)
	go func() {
		socket := WebSocket{}
		socket.ListenBack(func(bytes []byte) []byte {
			return bytes
		})
		socket.ListenAddr(Addr)
	}()
	go WsClient(channel)
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
	for {
		// 2 标识字节消息
		err := c.WriteMessage(2, []byte(message))
		if err != nil {
			break
		}
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("读取消息错误:", err)
			break
		}
		log.Println("获取数据:" + string(msg))
		channel <- 0
	}
}
