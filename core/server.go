package core

import "core/core"

type Server struct {
	Nacos core.Nacos
}

func (n *Server) Register(port uint64) {
	n.Nacos.SetServerConfig("192.168.101.10", 8848)
	n.Nacos.Register("192.168.101.10", port, core.ServerName)
}
