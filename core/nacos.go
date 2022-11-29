package core

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Nacos struct {
	registerParam vo.RegisterInstanceParam
	logoutParam   vo.DeregisterInstanceParam
	serverConfigs []constant.ServerConfig
	namingClient  naming_client.INamingClient
}

// Register 暂时先先写内部地址
func (n *Nacos) Register(ip string, port uint64) {
	port = 8848
	ip = "192.168.101.10"
	serviceName := "gateway-game"
	n.registerParam = vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	n.logoutParam = vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
	}

	// 创建serverConfig的另一种方式
	n.serverConfigs = []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}

	n.init()
}

func (n *Nacos) init() {
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{ServerConfigs: n.serverConfigs})
	if err != nil {
		print(err)
	}
	n.namingClient = namingClient
	success, err := n.namingClient.RegisterInstance(n.registerParam)
	if err != nil {
		print(err)
	}
	if success {
		fmt.Println("注册成功！")
	}
}

func (n *Nacos) Logout() {
	success, err := n.namingClient.DeregisterInstance(n.logoutParam)
	if err != nil {
		print(err)
	}
	if success {
		fmt.Println("注销成功！")
	}
}

func (n *Nacos) AllInstances() {
	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, err := n.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: n.registerParam.ServiceName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		print(err)
	}
	fmt.Println("所有实例: ", instances)
}
