package server

import "io-game-go/core"

type Server struct {
	Nacos core.Nacos
}

func (n *Server) Register() {
	n.Nacos.SetServerConfig("192.168.101.10", 8848)
	n.Nacos.Register("192.168.101.10", 8002, "server-game")
}
