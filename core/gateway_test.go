package core

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/messages"
	"log"
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
			ctx.Message.SetBody([]byte("用户"))
			login.SendAllUserMessage(ctx.Message.GetBytesResult())
		}
		ctx.Message.SetBody([]byte("广播"))
		gateway.SendMessage(ctx.Message.GetBytesResult())
		ctx.Message.SetBody([]byte("业务返回Hi->Ok->2"))
	})
	go gateway.Start(connect.Addr)
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

	clientInfo := registers.ClientInfo{Ip: "192.168.101.10", Port: 12344, //这里是rpc端口
		ServiceName: common.GatewayName, RemoteName: common.ServiceName} //测试时 RemoteName 传递一样的
	nacos := registers.NewNacos()
	nacos.RegisterClient(clientInfo)
	nacos.Register(registers.ServerInfo{Ip: "localhost", Port: 8848})

	gateway.SetClient(&remote.GrpcClient{})
	gateway.SetServer(&remote.GrpcServer{})
	gateway.SetRegisterClient(nacos)
	socket := &connect.WebSocket{}
	gateway.SetSocket(socket)
	gateway.Start(connect.Addr)
}

func WsTest(v chan int) {
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+connect.Addr+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(1, 1), Body: []byte("Hello, world!")}

	go func() {
		for true {
			// 读取消息
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			// 打印消息
			fmt.Printf("收到消息: %s\n", p)
			v <- 0
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
