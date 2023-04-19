package registers

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	// 创建两个etcd
	a1 := createTest(12345)
	a2 := createTest(12344)
	a1.Run()
	a2.Run()
	for i := 0; i < 6; i++ {
		ip, _ := a2.GetIp()
		fmt.Println(ip)
		list, _ := a2.ListIp(common.ServiceName)
		fmt.Println(list)
		time.Sleep(time.Second)
	}
	a1.Close()
	fmt.Println("a1 close")

	for i := 0; i < 3; i++ {
		ip, _ := a2.GetIp()
		fmt.Println(ip)
		time.Sleep(time.Second)
	}
	a2.Close()
	fmt.Println("a2 close")
}

func Test2(t *testing.T) {
	v := createTest(12300)
	v.Run()
	for i := 0; i < 10; i++ {
		ip, _ := v.GetIp()
		fmt.Println(ip)
		time.Sleep(time.Second)
	}
}

func createTest(port uint16) *Etcd {
	// 本地ip
	const localIp = "localhost"
	etcd := &Etcd{}
	etcd.RegisterServer(ServerInfo{Ip: localIp, Port: 2379})
	etcd.RegisterClient(ClientInfo{Ip: localIp, Port: port, ServiceName: common.ServiceName,
		RemoteName: common.ServiceName})
	return etcd
}
