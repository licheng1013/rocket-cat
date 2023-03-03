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
	//size := 0
	//v.uuidOnCoon.Range(func(key, value interface{}) bool {
	//	size++
	//	return true
	//})
	//common.FileLogger().Println("在线连接数:", size)

	v.AsyncResult(func(bytes []byte) {
		err = c.WriteMessage(2, bytes)
		if err != nil {
			// log.Println("写入错误:", err)
			common.FileLogger().Println("websocket写入错误: " + err.Error())
			_ = c.Close()
			v.uuidOnCoon.Delete(uuid)
		}
	})

	// 创建一个线程池，指定工作协程数为3，任务队列大小为10
	v.pool = common.NewPool(20, 30)
	v.pool.Start()

	for {
		// 1 字符串，2 字节
		_, message, err := c.ReadMessage()
		if err != nil {
			common.FileLogger().Println("websocket读取错误: " + err.Error())
			v.uuidOnCoon.Delete(uuid)
			break
		}
		v.InvokeMethod(message)

		//bytes := v.proxyMethod(message) // TODO 改成上面变为多线程进行发送数据
		//if len(bytes) == 0 {
		//	continue
		//}
		//err = c.WriteMessage(mt, bytes)
		//if err != nil {
		//	// log.Println("写入错误:", err)
		//	common.FileLogger().Println("websocket写入错误: " + err.Error())
		//	_ = c.Close()
		//	v.uuidOnCoon.Delete(uuid)
		//	break
		//}
	}

}
