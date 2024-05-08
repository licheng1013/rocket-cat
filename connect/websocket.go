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

func (ws *WebSocket) ListenBack(f func(uuid uint32, message []byte) []byte) {
	ws.ProxyMethod = f
}

func (ws *WebSocket) ListenAddr(addr string) {
	if ws.ProxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	// 判断Path是否为空并设置默认值/ws
	if ws.Path == "" {
		ws.Path = "/ws"
	}
	http.HandleFunc(ws.Path, ws.ws)
	if ws.Tls != nil {
		if err := http.ListenAndServeTLS(addr, ws.Tls.CertFile, ws.Tls.KeyFile, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}
}
func (ws *WebSocket) ws(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{
		// 排除前端跨域吹
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	uuid := ws.getNewChan()
	ws.AsyncResult(uuid, func(bytes []byte) {
		err = c.WriteMessage(websocket.BinaryMessage, bytes)
		if ws.handleErr(err, uuid, "websocket写入错误: ") {
			return
		}
	})

	for {
		// 1 字符串，2 字节
		_, message, err := c.ReadMessage()
		if ws.handleErr(err, uuid, "websocket读取错误: ") {
			break
		}
		ws.InvokeMethod(uuid, message)
	}

}
