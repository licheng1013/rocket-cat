package core

import (
	"github.com/io-game-go/registers"
	"github.com/io-game-go/remote"
	"log"
	"testing"
)

func TestService(t *testing.T) {
	rpc := &remote.GrpcServer{}
	rpc.SetRegister(registers.NewNacos())

	service := NewService()
	service.SetRpcServer(rpc)

	service.AddClose(rpc.Close)
	service.AddClose(func() {
		log.Println("在关机中了")
	})
	service.Start()
}
