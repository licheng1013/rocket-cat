package connect

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestWsServer(t *testing.T) {
	server(nil)
}

func TestWsTls(t *testing.T) {
	server(&Tls{KeyFile: "example.key", CertFile: "example.crt"})
}

func server(tls *Tls) {
	channel := make(chan int)
	socket := WebSocket{}
	socket.Tls = tls
	socket.OnClose(func(uuid uint32) {
		common.RocketLog.Println(uuid, "关闭了")
	})
	go func() {
		socket.ListenBack(func(uuid uint32, message []byte) []byte {
			common.RocketLog.Println(uuid)
			socket.SendMessage([]byte{}) // 测试空消息是否会返回
			//return message
			socket.SendSelectMessage([]byte("ok"), uuid)
			return []byte{}
		})
		socket.ListenAddr(Addr)
	}()

	// 客户端数量
	clientNum := 5
	for i := 0; i < clientNum; i++ {
		go WsClient(channel, tls != nil)
	}
	time.Sleep(1 * time.Second) // 需要等待1秒让客户端启动完成
	socket.SendMessage([]byte("广播消息-HelloWorld"))
	for i := 0; i < clientNum*2; i++ {
		select {
		case ok := <-channel:
			common.RocketLog.Println(ok)
		}
	}
}

func WsClient(channel chan int, enableTsl bool) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws"}
	if enableTsl {
		u.Scheme = "wss"
		websocket.DefaultDialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	go func() {
		for {
			_ = c.WriteMessage(websocket.BinaryMessage, []byte(HelloMsg)) //和广播一起测试
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
			common.RocketLog.Println("读取消息错误:", err)
			break
		}
		common.RocketLog.Println("获取数据:" + string(msg)) // 此处可能会打印多次,因为 channel <- 0 传输到通道也需要时间
		channel <- 0
	}
}
