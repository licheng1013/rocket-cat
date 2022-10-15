package main

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"io"
	"io-game-go/message"
	"io-game-go/router"
	"log"
	"testing"
)

func TestClient1(t *testing.T) {
	connectProto(1)
}

func TestClient2(t *testing.T) {
	connectProto(2)
}

// proto消息编码器
func connectProto(num int) {
	kecClient, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	info := message.Info{Info: "Hello" + fmt.Sprint(num)}
	marshal := message.MarshalBytes(&info)

	defaultMessage := message.ProtoMessage{Body: marshal, Merge: router.GetMerge(0, 1)}
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

			// TODO 这里是对数据处理实现部分，目前这个支持固定到字类

			v := message.Info{}
			err := proto.Unmarshal(buffer[:n], &v)
			if err != nil {
				log.Panicln(err)
			}
			log.Println("服务器消息: ", v.String())

		}
	}()
	for i := 0; i < 100; i++ {
		go func() {
			for true {
				marshal, err := proto.Marshal(&defaultMessage)
				if err != nil {
					log.Panicln(err)
				}
				_, _ = kecClient.Write(marshal)
			}
		}()
	}
	select {}
}

func connectJson(num int) {
	kecClient, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.DefaultMessage{Body: []byte("Hello" + fmt.Sprint(num)), Merge: router.GetMerge(0, 1)}
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

			// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
			log.Println(num, "服务端数据: ", string(buffer[:n]))
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
