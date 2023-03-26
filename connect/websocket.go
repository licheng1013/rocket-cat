// Package connect https://github.com/gorilla/websocket
package connect

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	MySocket
	// ws连接路径
	Path string
}

func (socket *WebSocket) ListenBack(f func(uuid uint32, message []byte) []byte) {
	socket.proxyMethod = f
}

func (socket *WebSocket) ListenAddr(addr string) {
	if socket.proxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	// 判断Path是否为空并设置默认值/ws
	if socket.Path == "" {
		socket.Path = "/ws"
	}
	http.HandleFunc(socket.Path, socket.ws)
	if socket.Tls != nil {
		if err := http.ListenAndServeTLS(addr, socket.Tls.CertFile, socket.Tls.KeyFile, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}
}
func (socket *WebSocket) ws(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{}
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	uuid := socket.getNewChan()
	socket.AsyncResult(uuid, func(bytes []byte) {
		err = c.WriteMessage(websocket.BinaryMessage, bytes)
		if socket.handleErr(err, uuid, "websocket写入错误: ") {
			return
		}
	})

	for {
		// 1 字符串，2 字节
		_, message, err := c.ReadMessage()
		if socket.handleErr(err, uuid, "websocket读取错误: ") {
			break
		}
		socket.InvokeMethod(uuid, message)
	}

}
