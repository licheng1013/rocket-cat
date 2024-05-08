package registers

import (
	"errors"
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"time"
)

// Nacos 请使用构造方法获取实例  NewNacos
type Nacos struct {
	namingClient  naming_client.INamingClient
	clientInfo    ClientInfo
	registerParam vo.RegisterInstanceParam
	logoutParam   vo.DeregisterInstanceParam
	updateParam   vo.UpdateInstanceParam
	serverInfo    ServerInfo
}

func (n *Nacos) Run() {
	// 创建serverConfig的另一种方式 -> 此处链接nacos的配置
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(n.serverInfo.Ip, uint64(n.serverInfo.Port), constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}
	var err error
	// 创建服务发现客户端的另一种方式 (推荐)
	n.namingClient, err = clients.NewNamingClient(vo.NacosClientParam{ServerConfigs: serverConfigs})
	if err != nil {
		panic(err)
	}
	fmt.Println(n.registerParam)
	if success, err := n.namingClient.RegisterInstance(n.registerParam); err != nil || !success {
		panic(err)
	}
	common.RocketLog.Println("注册中心:", n.serverInfo.Ip+":"+fmt.Sprint(n.serverInfo.Port))
	go n.heartbeat() // 心跳功能
}

func (n *Nacos) ClientInfo() ClientInfo {
	return n.clientInfo
}

func (n *Nacos) Close() {
	success, err := n.namingClient.DeregisterInstance(n.logoutParam)
	if err != nil {
		common.RocketLog.Println("注销错误:" + err.Error())
	}
	if success {
		common.RocketLog.Println("注销成功！")
	}
}

func (n *Nacos) RegisterServer(info ServerInfo) {
	n.serverInfo = info
}

// GetIp 获取单个ip
func (n *Nacos) GetIp() (ClientInfo, error) {
	// SelectList 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: n.clientInfo.RemoteName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		return ClientInfo{}, errors.New("获取实例为空")
	}
	return ClientInfo{Ip: instances.Ip, Port: uint16(instances.Port), ServiceName: instances.ServiceName}, nil
}

// ListIp 获取ip
func (n *Nacos) ListIp(serverName string) ([]ClientInfo, error) {
	instances := n.SelectList(serverName)
	infos := make([]ClientInfo, 0)
	if len(instances) == 0 {
		return infos, errors.New("获取实例为空")
	}
	for _, item := range instances {
		infos = append(infos, ClientInfo{Ip: item.Ip, Port: uint16(item.Port), ServiceName: item.ServiceName})
	}
	return infos, nil
}

func NewNacos() *Nacos {
	return &Nacos{}
}

// DefaultNacos 配置默认的nacos
func DefaultNacos() *Nacos {
	nacos := &Nacos{}
	clientInfo := ClientInfo{Ip: "localhost", Port: 12008,
		ServiceName: common.ServiceName, RemoteName: common.GatewayName}
	nacos.RegisterClient(clientInfo)
	nacos.RegisterServer(ServerInfo{Ip: "localhost", Port: 8848})
	return nacos
}

// TODO 这里作为保留教程
//func (n *Nacos) allInstances() {
//	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0 包括关闭得实例
//	instances, err := n.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
//		ServiceName: n.registerParam.ServiceName,
//		GroupName:   n.registerParam.GroupName,
//	})
//	if err != nil {
//		print(err)
//	}
//	common.Logger().Println("所有实例: ", instances)
//}

// SelectList 只获取存活的实例
func (n *Nacos) SelectList(serverName string) []model.Instance {
	// SelectList 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, _ := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
		HealthyOnly: true, //true 健康的实例
	})
	return instances
}

// Heartbeat 心跳功能
func (n *Nacos) heartbeat() {
	for true {
		instance, err := n.namingClient.UpdateInstance(n.updateParam)
		if err != nil || !instance {
			common.RocketLog.Println("更新实例失败,请检查Nacos!", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (n *Nacos) RegisterClient(info ClientInfo) {
	n.clientInfo = info
	// 这里是设置注册客户端的参数
	n.registerParam = vo.RegisterInstanceParam{
		Ip:          info.Ip,
		Port:        uint64(info.Port),
		ServiceName: info.ServiceName,
		GroupName:   "DEFAULT_GROUP", // 默认 default_GROUP
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	n.logoutParam = vo.DeregisterInstanceParam{
		Ip:          n.registerParam.Ip,
		Port:        n.registerParam.Port,
		ServiceName: n.registerParam.ServiceName,
		GroupName:   n.registerParam.GroupName,
	}
	n.updateParam = vo.UpdateInstanceParam{
		Ip:          n.registerParam.Ip,
		Port:        n.registerParam.Port,
		Weight:      n.registerParam.Weight,
		Enable:      n.registerParam.Enable,
		Healthy:     n.registerParam.Healthy,
		ServiceName: n.registerParam.ServiceName,
		GroupName:   n.registerParam.GroupName,
	}
}
