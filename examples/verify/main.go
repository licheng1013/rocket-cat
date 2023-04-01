package main

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"log"
)

func main() {
	var v int64
	socket := connect.KcpSocket{}
	socket.OnClose(func(socketId uint32) {
		// 关闭连接
		log.Println("关闭连接 -> ", socketId)
	})
	socket.Pool = common.NewPool()
	socket.ListenBack(func(uuid uint32, bytes []byte) []byte {
		fmt.Println("收到数据:" + string(bytes))
		v++
		return []byte("收到数据:" + string(bytes) + fmt.Sprint(v))
	})
	socket.ListenAddr("127.0.0.1:12355")
}
