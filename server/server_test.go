package server

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	gateway := &Server{}
	gateway.Register()
	time.Sleep(5 * time.Second)
	gateway.Nacos.AllInstances()
	time.Sleep(30 * time.Second)
	gateway.Nacos.Logout()
}
