package server

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server1 := test(8002)
	server2 := test(8003)
	time.Sleep(100 * time.Second)
	server1.Nacos.Logout()
	server2.Nacos.Logout()
}

func test(port uint64) *Server {
	server := &Server{}
	server.Register(port)
	return server
}
