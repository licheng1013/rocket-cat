package core

import (
	"core/connect"
	"testing"
)

func TestGateway(t *testing.T) {
	gateway := NewGateway()
	gateway.Start(connect.Addr, &connect.WebSocket{})
}
