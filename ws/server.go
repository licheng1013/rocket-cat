package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// Server 封装 WebSocket 升级、收包和路由分发流程。
type Server struct {
	router   *Router            // 路由器
	upgrader websocket.Upgrader // WebSocket 升级器

	nextSessionId int64 // 下一个会话 ID
}

// NewServer 创建一个 WebSocket 服务入口。
func NewServer(router *Router) *Server {
	return &Server{
		router: router,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// ServeHTTP 处理 WebSocket 连接并循环分发消息。
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	session := NewSession(atomic.AddInt64(&s.nextSessionId, 1), conn)
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var packet Packet
		if err := json.Unmarshal(payload, &packet); err != nil {
			Fail(&Context{Session: session}, 400, err.Error())
			continue
		}

		s.router.Dispatch(&Context{
			Session: session,
			Packet:  &packet,
		})
	}
}
