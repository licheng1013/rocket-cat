package gateway

import (
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	gateway := &Gateway{}
	gateway.Register()
	time.Sleep(5 * time.Second)
	gateway.Nacos.AllInstances()
	time.Sleep(30 * time.Second)
	gateway.Nacos.Logout()
}
