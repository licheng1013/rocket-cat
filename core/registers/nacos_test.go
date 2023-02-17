package registers

import (
	"fmt"
	"github.com/io-game-go/common"
	"testing"
	"time"
)

func TestNacos(t *testing.T) {
	nacos := NewNacos()
	// 这里一致用于测试访问
	info := RegisterInfo{Ip: "192.168.101.10", Port: 8848, ServiceName: common.ServicerName, RemoteName: common.ServicerName}
	nacos.Register(info)
	time.Sleep(3 * time.Second)
	fmt.Println(nacos.ListIp())
	time.Sleep(10 * time.Second)
	nacos.Close()
}
