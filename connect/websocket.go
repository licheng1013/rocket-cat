// Package connect https://github.com/gorilla/websocket
package connect

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
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
	socket.init()
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
	uuid := common.UuidKit.UUID()
	messageChannel := make(chan []byte)
	socket.UuidOnCoon.Store(uuid, messageChannel)

	// 统计数
	//size := 0
	//socket.UuidOnCoon.Range(func(key, value interface{}) bool {
	//	size++
	//	return true
	//})
	//common.Logger().Println("在线连接数:", size)

	socket.AsyncResult(uuid, func(bytes []byte) {
		err = c.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			// router.FileLogger().Println("写入错误:", err)
			common.FileLogger().Println("websocket写入错误: " + err.Error())
			_ = c.Close()
			socket.close(uuid)
		}
	})

	for {
		// 1 字符串，2 字节
		_, message, err := c.ReadMessage()
		if err != nil {
			common.FileLogger().Println("websocket读取错误: " + err.Error())
			socket.close(uuid)
			break
		}
		socket.InvokeMethod(uuid, message)

		//bytes := socket.proxyMethod(message) // TODO 改成上面变为多线程进行发送数据
		//if len(bytes) == 0 {
		//	continue
		//}
		//err = c.WriteMessage(mt, bytes)
		//if err != nil {
		//	// router.FileLogger().Println("写入错误:", err)
		//	common.Logger().Println("websocket写入错误: " + err.Error())
		//	_ = c.Close()
		//	socket.UuidOnCoon.Delete(uuid)
		//	break
		//}
	}

}
