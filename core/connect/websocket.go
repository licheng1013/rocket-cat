// Package connect https://github.com/gorilla/websocket
package connect

import (
	"github.com/gorilla/websocket"
	"github.com/io-game-go/common"
	"net/http"
)

type WebSocket struct {
	MySocket
}

func (v *WebSocket) ListenBack(f func([]byte) []byte) {
	v.proxyMethod = f
}

func (v *WebSocket) ListenAddr(addr string) {
	if v.proxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	http.HandleFunc("/ws", v.ws)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
func (v *WebSocket) ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	uuid := common.UuidKit.UUID()
	v.uuidOnCoon.Store(uuid, c)
	// 统计数
	size := 0
	v.uuidOnCoon.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	common.FileLogger().Println("在线连接数:", size)
	for {
		// 1 字符串，2 字节
		mt, message, err := c.ReadMessage()
		if err != nil {
			common.FileLogger().Println("websocket读取错误: " + err.Error())
			v.uuidOnCoon.Delete(uuid)
			break
		}
		// log.Printf("收到消息: %s", message)
		bytes := v.proxyMethod(message)
		if len(bytes) == 0 {
			continue
		}
		err = c.WriteMessage(mt, bytes)
		if err != nil {
			// log.Println("写入错误:", err)
			common.FileLogger().Println("websocket写入错误: " + err.Error())
			_ = c.Close()
			v.uuidOnCoon.Delete(uuid)
			break
		}
	}
}
