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
	info := RegisterInfo{Ip: "192.168.101.10", Port: 8848}
	clientInfo := RegisterInfo{Ip: "192.168.101.10", Port: 12345,
		ServiceName: common.ServicerName, RemoteName: common.ServicerName} // 测试时 RemoteName 传递一样的
	nacos.RegisterClient(clientInfo)
	nacos.Register(info)
	time.Sleep(5 * time.Second)
	fmt.Println(nacos.GetIp())
	time.Sleep(10 * time.Second)
	nacos.Close()
}
