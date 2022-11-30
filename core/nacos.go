package core

import (
	"fmt"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io"
	"log"
	"strconv"
	"time"
)

type Nacos struct {
	registerParam vo.RegisterInstanceParam
	logoutParam   vo.DeregisterInstanceParam
	serverConfigs []constant.ServerConfig
	namingClient  naming_client.INamingClient
}

func (n *Nacos) SetServerConfig(ip string, port uint64) {
	// 创建serverConfig的另一种方式 -> 此处链接nacos的配置
	n.serverConfigs = []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}
}

// Register 注册进入Nacos,必须先调用 SetServerConfig 才能使用
// 这里设置的都是客户端地址和ip远程调用获取则是这里注册的 ip 和 端口
func (n *Nacos) Register(ip string, port uint64, serviceName string) {
	if len(n.serverConfigs) == 0 {
		print("未设置Nacos服务端配置")
	}
	n.registerParam = vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	n.logoutParam = vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
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
		n.heartbeat()
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

// SelectInstances 获取指定条件的实例
func (n *Nacos) SelectInstances(serverName string) []model.Instance {
	// SelectInstances 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
		HealthyOnly: true, //true 健康的实例
	})
	if err != nil {
		print(err)
	}
	return instances
}

// SelectOneHealthyInstance 获取指定条件的实例
func (n *Nacos) SelectOneHealthyInstance(serverName string) *model.Instance {
	// SelectInstances 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		print(err)
	}
	return instances
}

// 心跳功能
func (n *Nacos) heartbeat() {
	go func() {
		for true {
			time.Sleep(1 * time.Second)
			url := "/nacos/v2/ns/health/instance"
			heads := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
			params := map[string]string{}
			params["serviceName"] = n.registerParam.ServiceName
			params["ip"] = n.registerParam.Ip
			params["port"] = strconv.FormatUint(n.registerParam.Port, 10)
			params["healthy"] = "true"

			for _, config := range n.serverConfigs {
				httpUrl := fmt.Sprintf("http://%v:%v%v", config.IpAddr, strconv.FormatUint(config.Port, 10), url)
				fmt.Println("心跳地址: " + httpUrl)
				resp, err := netutil.HttpPut(httpUrl, heads, params)
				if err != nil {
					log.Println(err)
				}
				body, _ := io.ReadAll(resp.Body)
				fmt.Println("心跳结果: " + string(body))
			}
		}
	}()
}
