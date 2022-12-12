package main

import (
	"core/message"
	"core/router"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"time"
)

// 测试客户端连接
func main() {
	connectJson(1)
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
		for true {
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
	go func() {
		for true {
			marshal, err := proto.Marshal(&defaultMessage)
			if err != nil {
				log.Panicln(err)
			}
			_, _ = kecClient.Write(marshal)
		}
	}()
	select {}
}

func connectJson(num int) {
	kecClient, err := kcp.DialWithOptions("localhost:8001", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.JsonMessage{Body: []byte("Hello" + fmt.Sprint(num)), Merge: router.GetMerge(0, 1)}
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
	go func() {
		for true {
			_, _ = kecClient.Write(message.GetObjectToBytes(defaultMessage))
			time.Sleep(1 * time.Second)
		}
	}()
	select {}
}
