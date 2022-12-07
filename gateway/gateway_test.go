package gateway

import (
	"fmt"
	"io-game-go/core"
	"testing"
)

func TestGateway(t *testing.T) {
	gateway := &Gateway{}
	gateway.Register()
	instances1 := gateway.Nacos.SelectOneHealthyInstance(core.ServerName)
	instances2 := gateway.Nacos.SelectOneHealthyInstance(core.ServerName)
	fmt.Println(instances1.Ip)
	fmt.Println(instances2.Ip)
	gateway.Nacos.Logout()
}
