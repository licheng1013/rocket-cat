package core

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
	"github.com/licheng1013/rocket-cat/router"
)

func TestSingleGateway(t *testing.T) {
	channel := make(chan int)
	gateway := DefaultGateway()
	gateway.Router().AddAction(1, 1, func(ctx *router.Context) {
		r := gateway.GetPlugin(LoginPluginId)
		login := r.(LoginInterface)
		if login.Login(12345, ctx.SocketId) {
			fmt.Printf("login.ListUserId(): %v\n", login.ListUserId())
			login.SendAllUserMessage(ctx.Message.SetBody([]byte("用户")).GetBytesResult())
			gateway.SendMessage(ctx.Message.SetBody([]byte("广播")).GetBytesResult())
		}

		ctx.Message.SetBody([]byte("业务返回Hi->Ok->2"))
	})
	go gateway.Start(connect.Addr, &connect.WebSocket{})
	time.Sleep(time.Second * 1) //等待完全启动
	go WsTest(channel)
	select {
	case ok := <-channel:
		time.Sleep(time.Second * 1)
		common.Logger().Println(ok)
	}
}

func TestGateway(t *testing.T) {
	gateway := NewGateway()
	gateway.SetSingle(false)

	clientInfo := registers.RegisterInfo{Ip: "192.168.101.10", Port: 12344, //这里是rpc端口
		ServiceName: common.GatewayName, RemoteName: common.ServiceName} //测试时 RemoteName 传递一样的
	nacos := registers.NewNacos()
	nacos.RegisterClient(clientInfo)
	nacos.Register(registers.RegisterInfo{Ip: "localhost", Port: 8848})

	gateway.SetClient(&remote.GrpcClient{})
	gateway.SetServer(&remote.GrpcServer{})
	gateway.SetRegisterClient(nacos)
	gateway.Start(connect.Addr, &connect.WebSocket{})
}

func WsTest(v chan int) {
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
	var count int64
	go func() {
		defer close(done)
		for {
			_, m, err := c.ReadMessage()
			jsonDecoder := decoder.JsonDecoder{}
			dto := jsonDecoder.DecoderBytes(m)
			if err != nil {
				common.Logger().Println("读取消息错误:", err)
				return
			}
			log.Printf("收到消息-> %v", string(dto.GetBody()))
			count++
			if count >= 3 {
				v <- 0
			}
		}
	}()
	var b bool
	for {
		if b {
			continue
		}
		jsonMessage := messages.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.Merge = common.CmdKit.GetMerge(1, 1)
		err := c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
		if err != nil {
			common.Logger().Println("写:", err)
			return
		}
		b = true
	}
}
