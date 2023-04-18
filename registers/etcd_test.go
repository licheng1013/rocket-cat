package registers

import (
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
	for i := 0; i < 5; i++ {
		_, _ = a2.GetIp()
		time.Sleep(time.Second)
	}
	a1.Close()
	for i := 0; i < 5; i++ {
		_, _ = a2.GetIp()
		time.Sleep(time.Second)
	}
	a2.Close()
}

func createTest(port uint16) *Etcd {
	// 本地ip
	const localIp = "localhost"
	etcd := &Etcd{}
	etcd.ServerInfo = ServerInfo{Ip: localIp, Port: 2379}
	etcd.ClientInfo = ClientInfo{Ip: localIp, Port: port, ServiceName: common.ServiceName,
		RemoteName: common.ServiceName}
	return etcd
}
