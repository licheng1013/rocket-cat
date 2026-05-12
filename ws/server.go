package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type Server struct {
	router   *Router
	upgrader websocket.Upgrader

	nextSessionId int64
}

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
