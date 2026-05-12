package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Session struct {
	Id   int64
	Conn *websocket.Conn

	Uid int64

	Data map[string]any

	writeMu sync.Mutex
}

func NewSession(id int64, conn *websocket.Conn) *Session {
	return &Session{
		Id:   id,
		Conn: conn,
		Data: make(map[string]any),
	}
}

func (s *Session) SendJSON(v any) error {
	if s == nil || s.Conn == nil {
		return nil
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	return s.Conn.WriteJSON(v)
}
