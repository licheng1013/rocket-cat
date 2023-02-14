// Package connect https://github.com/gorilla/websocket
package connect

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocket struct {
	funcMsg func([]byte) []byte
}

func (v *WebSocket) ListenBack(f func([]byte) []byte) {
	v.funcMsg = f
}

func (v *WebSocket) ListenAddr(addr string) {
	if v.funcMsg == nil {
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
		log.Println("升级为WebSocket错误:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("读取错误:", err)
			break
		}
		// log.Printf("收到消息: %s", message)
		bytes := v.funcMsg(message)
		if len(bytes) == 0 {
			continue
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("写入错误:", err)
			break
		}
	}
}
