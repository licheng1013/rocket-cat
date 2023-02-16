package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/io-game-go/connect"
	"github.com/io-game-go/decoder"
	"github.com/io-game-go/message"
	"github.com/io-game-go/router"
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

	start := time.Now().UnixMilli()
	var count int64
	gateway.Router().AddFunc(10, func(ctx router.Context) []byte {
		ctx.Message.SetBody([]byte("Hi Ok"))
		end := time.Now().UnixMilli()
		count++
		if end-start > 1000 {
			fmt.Println("1s请求数量:", count)
			count = 0
			start = end
		}
		//log.Println(string(ctx.Message.GetBody()))
		return ctx.Message.GetBytesResult()
	})
	fmt.Println(start)
	gateway.Start(connect.Addr, &connect.WebSocket{})
}

func TestWsClient2(t *testing.T) {
	for i := 0; i < 2; i++ {
		go WsTest()
	}
	WsTest()
}

func WsTest() {
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
	for {
		jsonMessage := message.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.Merge = 10
		err := c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
		if err != nil {
			log.Println("写:", err)
			return
		}

	}
}
