package main

import (
	"io-game-go/core"
	"io-game-go/router"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	// 默认的消息实现: DefaultMessage
	router.AddFunc(router.GetMerge(0, 1), func(msg interface{}) {
		log.Println("收到消息: ", msg)
		log.Println("回调成功！")
	})

	server := core.NewGameServer()
	server.Run()
}
