package connect

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func ws(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("收到消息: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("写入错误:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", ws)
	if err := http.ListenAndServe(Addr, nil); err != nil {
		panic(err)
	}
}
