package core

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/remote"
	"github.com/licheng1013/rocket-cat/router"
)

func TestSingleGateway(t *testing.T) {
	channel := make(chan int)
	gateway := DefaultGateway()
	gateway.Router().AddAction(1, 1, func(ctx *router.Context) {
		log.Println("获取数据 -> ", string(ctx.Message.GetBody()))
		r := gateway.GetPlugin(LoginPluginId)
		login := r.(LoginInterface)
		if login.Login(12345, ctx.SocketId) {
			fmt.Printf("login.ListUserId(): %v\n", login.ListUserId())
			login.SendAllUserMessage(ctx.Message.SetBody([]byte("用户")).GetBytesResult())
		}
		gateway.SendMessage(ctx.Message.SetBody([]byte("广播")).GetBytesResult())
		ctx.Message.SetBody([]byte("业务返回Hi->Ok->2"))
	})
	socket := &connect.WebSocket{}
	socket.Debug = true
	go gateway.Start(connect.Addr, socket)
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
	u := url.URL{Scheme: "ws", Host: connect.Addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	for {
		// 只写一次就可以了
		jsonMessage := messages.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.Merge = common.CmdKit.GetMerge(1, 1)
		err = c.WriteMessage(websocket.BinaryMessage, jsonMessage.GetBytesResult())
		if err != nil {
			common.Logger().Println("写:", err)
		}
		_, m, err := c.ReadMessage()
		jsonDecoder := decoder.JsonDecoder{}
		dto := jsonDecoder.DecoderBytes(m)
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("收到消息-> %v", string(dto.GetBody()))
		time.Sleep(time.Second)
	}

}
