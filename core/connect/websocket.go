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

func (socket *WebSocket) ListenBack(f func([]byte) []byte) {
	socket.proxyMethod = f
}

func (socket *WebSocket) ListenAddr(addr string) {
	if socket.proxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	http.HandleFunc("/ws", socket.ws)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
func (socket *WebSocket) ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	uuid := common.UuidKit.UUID()
	socket.uuidOnCoon.Store(uuid, c)
	// 统计数
	//size := 0
	//socket.uuidOnCoon.Range(func(key, value interface{}) bool {
	//	size++
	//	return true
	//})
	//common.FileLogger().Println("在线连接数:", size)

	socket.AsyncResult(func(bytes []byte) {
		err = c.WriteMessage(2, bytes)
		if err != nil {
			// log.Println("写入错误:", err)
			common.FileLogger().Println("websocket写入错误: " + err.Error())
			_ = c.Close()
			socket.uuidOnCoon.Delete(uuid)
		}
	})

	socket.init()

	for {
		// 1 字符串，2 字节
		_, message, err := c.ReadMessage()
		if err != nil {
			common.FileLogger().Println("websocket读取错误: " + err.Error())
			socket.uuidOnCoon.Delete(uuid)
			break
		}
		socket.InvokeMethod(message)

		//bytes := socket.proxyMethod(message) // TODO 改成上面变为多线程进行发送数据
		//if len(bytes) == 0 {
		//	continue
		//}
		//err = c.WriteMessage(mt, bytes)
		//if err != nil {
		//	// log.Println("写入错误:", err)
		//	common.FileLogger().Println("websocket写入错误: " + err.Error())
		//	_ = c.Close()
		//	socket.uuidOnCoon.Delete(uuid)
		//	break
		//}
	}

}
