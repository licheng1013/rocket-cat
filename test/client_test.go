package main

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"io-game-go/message"
	"io-game-go/router"
	"log"
	"testing"
)

func TestClient1(t *testing.T) {
	connect(1)
}

func TestClient2(t *testing.T) {
	connect(2)
}

func connect(num int) {
	kecClient, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.DefaultMessage{Body: "Hello" + fmt.Sprint(num), Merge: router.GetMerge(0, 1)}
	// 获取服务单的消息
	go func() {
		var buffer = make([]byte, 1024, 1024)
		for {
			// 读取长度 n
			n, e := kecClient.Read(buffer)
			if e != nil {
				if e == io.EOF {
					break
				}
				fmt.Println(errorx.Wrap(e))
				break
			}

			msg := message.GetBytesToObject(buffer[:n])
			// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
			m := message.DefaultMessage{}
			router.GetObjectToToMap(msg, &m)
			log.Println(num, "服务端数据: ", m)
		}
	}()
	for i := 0; i < 100; i++ {
		go func() {
			for true {
				_, _ = kecClient.Write(message.GetObjectToBytes(defaultMessage))
			}
		}()
	}
	select {}
}
