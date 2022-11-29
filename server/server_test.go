package server

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server1 := test(8001)
	time.Sleep(100 * time.Second)
	server1.Nacos.Logout()
}

func test(port uint64) *Server {
	server := &Server{}
	server.Register(port)
	return server
}
