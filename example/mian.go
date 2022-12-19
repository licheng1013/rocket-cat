package main

import (
	"core/message"
	"core/protof"
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
	for i := 0; i < 2; i++ {
		go func(i int) {
			connectProto(i)
		}(i)
	}
	select {}
}

// proto消息编码器
func connectProto(num int) {
	kecClient, err := kcp.DialWithOptions("localhost:8001", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	info := "HelloWorld"
	defaultMessage := protof.ProtoMessage{Body: []byte(info), Merge: router.GetMerge(1, 1)}
	// 获取服务单的消息
	unix := time.Now().UnixMilli()
	var count int64
	fmt.Println(unix)
	// 获取服务单的消息
	go func() {
		var buffer = make([]byte, 1024*8)
		for true {
			count++
			// 读取长度 n
			n, e := kecClient.Read(buffer)
			if e != nil {
				if e == io.EOF {
					continue
				}
				fmt.Println(errorx.Wrap(e))
				continue
			}

			// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
			v := protof.ProtoMessage{}
			err := proto.Unmarshal(buffer[:n], &v)
			if err != nil {
				log.Panicln(err)
			}
			// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
			log.Println(num, "服务端数据: ", string(v.Body))

			newUnix := time.Now().UnixMilli()
			if newUnix-unix > 1000 {
				fmt.Println(fmt.Sprint(num)+"线程1秒请求数:", fmt.Sprint(count))
				unix = newUnix
				count = 0
			}

		}
	}()
	go func() {
		for true {
			marshal, err := proto.Marshal(&defaultMessage)
			if err != nil {
				log.Panicln(err)
			}
			_, err = kecClient.Write(marshal)
			//log.Println("写入: "+fmt.Sprint(num))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	select {}
}

func connectJson(num int) {
	kecClient, err := kcp.DialWithOptions("localhost:8001", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.JsonMessage{Body: []byte("Hello" + fmt.Sprint(num)), Merge: router.GetMerge(1, 1)}
	// 获取服务单的消息
	unix := time.Now().UnixMilli()
	fmt.Println(unix)
	var count int64
	go func() {
		var buffer = make([]byte, 1024, 1024)
		for {
			count++
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

			newUnix := time.Now().UnixMilli()
			if newUnix-unix > 1000 {
				fmt.Println("1秒请求数:", fmt.Sprint(count))
				unix = newUnix
				count = 0
			}
		}
	}()
	go func() {
		for true {
			_, _ = kecClient.Write(message.GetObjectToBytes(defaultMessage))
			//time.Sleep(1 * time.Second)
		}
	}()
	select {}
}
