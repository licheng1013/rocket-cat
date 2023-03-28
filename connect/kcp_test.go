package connect

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"testing"
	"time"
)

func TestKcpServer(t *testing.T) {
	channel := make(chan int)
	socket := KcpSocket{}
	socket.Pool = common.NewPool()
	socket.ListenBack(func(uuid uint32, bytes []byte) []byte {
		return bytes
	})
	go socket.ListenAddr("localhost:12355")
	time.Sleep(time.Second)
	go KcpClient(channel)
	select {
	case ok := <-channel:
		common.Logger().Println(ok)
	}
}

func KcpClient(channel chan int) {
	//common.Logger().Println("客户端监听:" + Addr)
	if client, err := kcp.DialWithOptions("localhost:12355", nil, 10, 3); err == nil {
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
				common.Logger().Println("获取数据:" + string(buf))
				channel <- 0
			} else {
				log.Fatal(err)
			}
		}
	} else {
		common.Logger().Println("监听异常:", err)
	}
}
