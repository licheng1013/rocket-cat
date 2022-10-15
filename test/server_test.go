package main

import (
	"google.golang.org/protobuf/proto"
	"io-game-go/core"
	"io-game-go/message"
	"io-game-go/router"
	"log"
	"testing"
)

func TestJsonServer(t *testing.T) {
	// 默认的消息实现: DefaultMessage
	router.AddFunc(router.GetMerge(0, 1), func(msg interface{}) interface{} {
		log.Println("收到消息: ", string(msg.([]byte)))
		return msg
	})

	server := core.NewGameServer()
	server.Run()
}

func TestProtoServer(t *testing.T) {
	// 默认的消息实现: DefaultMessage
	router.AddFunc(router.GetMerge(0, 1), func(msg interface{}) interface{} {

		info := message.Info{}
		// 转换反序列话
		err := proto.Unmarshal(msg.([]byte), &info)
		if err != nil {
			panic(err)
		}
		log.Println(info.String())

		return msg
	})

	server := core.NewGameServer()
	server.Run()
}
