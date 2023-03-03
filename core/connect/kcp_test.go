package connect

import (
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"testing"
)

func TestKcpServer(t *testing.T) {
	channel := make(chan int)
	socket := KcpSocket{}
	go func() {
		socket.ListenBack(func(bytes []byte) []byte {
			return bytes
		})
		socket.ListenAddr(Addr)
	}()
	go KcpClient(channel)
	select {
	case ok := <-channel:
		log.Println(ok)
	}
}

func KcpClient(channel chan int) {
	//log.Println("客户端监听:" + Addr)
	if client, err := kcp.DialWithOptions(Addr, nil, 10, 3); err == nil {
		for {
			data := HelloMsg
			buf := make([]byte, len(data))
			if _, err := client.Write([]byte(data)); err == nil {
				if _, err := io.ReadFull(client, buf); err == nil {
					log.Println("获取数据:" + string(buf))
					channel <- 0
				} else {
					log.Fatal(err)
				}
			}
		}
	} else {
		log.Println("监听异常:", err)
	}
}
