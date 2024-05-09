package connect

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"runtime"
	"testing"
	"time"
)

func TestKcpServer(t *testing.T) {
	go func() {
		// 打印线程数
		for true {
			fmt.Println("协程数 -> ", runtime.NumGoroutine())
			time.Sleep(3 * time.Second)
		}
	}()
	socket := KcpSocket{}
	socket.OnClose(func(socketId uint32) {
		// 关闭连接
		log.Println("关闭连接 -> ", socketId)
	})
	socket.Pool = common.NewPool()
	socket.ListenBack(func(uuid uint32, bytes []byte) []byte {
		fmt.Println("收到数据:" + string(bytes))
		return bytes
	})
	//socket.ListenAddr("127.0.0.1:12355")
	go socket.ListenAddr("localhost:12355")
	channel := make(chan int)
	time.Sleep(time.Second)
	go KcpClient(channel)
	select {
	case ok := <-channel:
		common.CatLog.Println(ok)
	}
}

func KcpClient(channel chan int) {
	//common.Logger().Println("客户端监听:" + Addr)
	if client, err := kcp.DialWithOptions("localhost:12355", nil, 0, 0); err == nil {
		data := HelloMsg
		buf := make([]byte, len(data))
		go func() {
			for {
				fmt.Println("发送数据:" + data)
				_, _ = client.Write([]byte(data))
				time.Sleep(time.Second * 2)
			}
		}()
		for {
			if _, err := io.ReadFull(client, buf); err == nil {
				common.CatLog.Println("获取数据:" + string(buf))
				channel <- 0
			} else {
				log.Fatal(err)
			}
		}
	} else {
		common.CatLog.Println("监听异常:", err)
	}
}
