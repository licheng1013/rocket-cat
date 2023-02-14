package connect

import (
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"testing"
	"time"
)

func TestKcpServer(t *testing.T) {
	socket := KcpSocket{}
	socket.ListenBack(func(bytes []byte) []byte {
		log.Println("收到消息:" + string(bytes))
		return bytes
	})
	socket.ListenAddr(Addr)
}

func TestKcpClient(t *testing.T) {
	log.Println("客户端监听:" + Addr)
	if client, err := kcp.DialWithOptions(Addr, nil, 10, 3); err == nil {
		for {
			buf := make([]byte, len(HelloMsg))
			if _, err := client.Write([]byte(HelloMsg)); err == nil {
				if _, err := io.ReadFull(client, buf); err == nil {
					log.Println("返回:", string(buf))
				} else {
					log.Fatal(err)
				}
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Println("监听异常:", err)
	}
}
