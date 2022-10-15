package main

import (
	"io-game-go/core"
	"io-game-go/decoder"
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
		message.UnmarshalInterface(msg, &info)
		log.Println(info.String())
		return msg
	})

	server := core.NewGameServer()
	// 设置编码器
	server.SetDecoder(decoder.ProtoDecoder{})
	server.Run()
}
