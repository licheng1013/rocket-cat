package gateway

import (
	"fmt"
	"io-game-go/core"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	gateway := &Gateway{}
	gateway.Register()
	instances1 := gateway.Nacos.SelectOneHealthyInstance(core.ServerName)
	instances2 := gateway.Nacos.SelectOneHealthyInstance(core.ServerName)
	fmt.Println(instances1)
	fmt.Println(instances2)
	time.Sleep(30 * time.Second)
	gateway.Nacos.Logout()
}
