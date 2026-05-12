package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Session 表示一个客户端连接会话。
type Session struct {
	Id   int64           // 会话 ID
	Conn *websocket.Conn // WebSocket 连接

	Uid int64 // 绑定的用户 ID

	Data map[string]any // 会话级临时数据

	writeMu sync.Mutex
}

// NewSession 创建一个客户端连接会话。
func NewSession(id int64, conn *websocket.Conn) *Session {
	return &Session{
		Id:   id,
		Conn: conn,
		Data: make(map[string]any),
	}
}

// SendJSON 向当前会话写入 JSON 响应。
func (s *Session) SendJSON(v any) error {
	if s == nil || s.Conn == nil {
		return nil
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	return s.Conn.WriteJSON(v)
}
