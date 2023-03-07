package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/licheng1013/io-game-go/common"
	"github.com/licheng1013/io-game-go/connect"
	"github.com/licheng1013/io-game-go/decoder"
	"github.com/licheng1013/io-game-go/messages"
	"github.com/licheng1013/io-game-go/registers"
	"github.com/licheng1013/io-game-go/remote"
	"github.com/licheng1013/io-game-go/router"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestSingleGateway(t *testing.T) {
	socket := &connect.WebSocket{}
	channel := make(chan int)

	gateway := NewGateway()
	gateway.SetDecoder(decoder.JsonDecoder{})
	start := time.Now().UnixMilli()
	var count int64
	gateway.Router().AddAction(common.CmdKit.GetMerge(1, 1), func(ctx *router.Context) {
		end := time.Now().UnixMilli()
		count++
		socket.SendMessage(ctx.Message.SetBody([]byte("Hi")).GetBytesResult())
		ctx.Message.SetBody([]byte("Hi Ok 2"))
		if end-start > 100 { //此处设置监听请求时间 -> 配置和Idea: RequestTool 插件能够调试广播，或者在开个客户端
			fmt.Println("1s请求数量:", count)
			count = 0
			start = end
			channel <- 0
		}
		//log.Println(string(ctx.Message.GetBody()))
		//ctx.Message = nil
	})
	fmt.Println(start)

	go gateway.Start(connect.Addr, socket)
	go WsTest()
	select {
	case ok := <-channel:
		log.Println(ok)
	}

}

func TestGateway(t *testing.T) {

	gateway := NewGateway()
	gateway.SetSingle(false)

	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: 12345,
		ServiceName: common.GatewayName, RemoteName: common.ServicerName} // 测试时 RemoteName 传递一样的
	nacos := registers.NewNacos()
	nacos.RegisterClient(clientInfo)
	nacos.Register(registers.RegisterInfo{Ip: "localhost", Port: 8848})

	gateway.SetClient(&remote.GrpcClient{})
	gateway.SetRegisterClient(nacos)
	gateway.Start(connect.Addr, &connect.WebSocket{})
}

func TestWsClient2(t *testing.T) {
	//for i := 0; i < 2; i++ {
	//	go WsTest()
	//}
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
		jsonMessage := messages.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.Merge = common.CmdKit.GetMerge(1, 1)
		err := c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
		if err != nil {
			log.Println("写:", err)
			return
		}

	}
}
