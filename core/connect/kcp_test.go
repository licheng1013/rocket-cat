package connect

import (
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	ListenerKcp(Addr)
}

func TestClient(t *testing.T) {
	log.Println("客户端监听:" + Addr)
	if client, err := kcp.DialWithOptions(Addr, nil, 10, 3); err == nil {
		for {
			msg := "HelloWorld"
			buf := make([]byte, len(msg))
			if _, err := client.Write([]byte(msg)); err == nil {
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
