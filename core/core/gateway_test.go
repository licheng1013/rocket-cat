package core

import (
	"core/connect"
	"core/decoder"
	"core/message"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	gateway := NewGateway()
	gateway.SetDecoder(decoder.JsonDecoder{})

	gateway.Router().AddFunc(10, func(msg message.Message) []byte {
		log.Println(string(msg.GetBody()))
		return msg.GetBody()
	})

	gateway.Start(connect.Addr, &connect.WebSocket{})
}

func TestWsClient(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: connect.Addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, m, err := c.ReadMessage()
			jsonDecoder := decoder.JsonDecoder{}
			dto := jsonDecoder.DecoderBytes(m)
			if err != nil {
				log.Println("读取消息错误:", err)
				return
			}
			log.Println("收到消息:", string(dto.GetBody()))
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			jsonMessage := message.JsonMessage{Body: []byte(t.String())}
			jsonMessage.Merge = 10
			err := c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
			if err != nil {
				log.Println("写:", err)
				return
			}
		case <-interrupt:
			log.Println("中断")
			// 通过发送关闭消息干净地关闭连接，然后等待（超时）服务器关闭连接。
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("写关闭:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
