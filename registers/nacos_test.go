package registers

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"testing"
	"time"
)

func TestNacos(t *testing.T) {
	const localIp = "localhost"
	nacos := NewNacos()
	// 客户端的注册信息
	clientInfo := ClientInfo{Ip: localIp, Port: 12345,
		ServiceName: common.ServiceName, RemoteName: common.ServiceName} // 测试时 RemoteName 传递一样的
	nacos.RegisterClient(clientInfo)
	// nacos的注册地址
	nacos.RegisterServer(ServerInfo{Ip: localIp, Port: 8848})
	time.Sleep(5 * time.Second)
	fmt.Println(nacos.GetIp())
	time.Sleep(10 * time.Second)
	nacos.Close()
}
