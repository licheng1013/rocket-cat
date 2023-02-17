package registers

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"time"
)

// Nacos 请使用构造方法获取实例  NewNacos
type Nacos struct {
	registerParam vo.RegisterInstanceParam
	logoutParam   vo.DeregisterInstanceParam
	serverConfigs []constant.ServerConfig
	namingClient  naming_client.INamingClient
	registerInfo  RegisterInfo
}

func (n *Nacos) Close() {
	success, err := n.namingClient.DeregisterInstance(n.logoutParam)
	if err != nil {
		print(err)
	}
	if success {
		log.Println("注销成功！")
	}
}

func (n *Nacos) Register(info RegisterInfo) {
	n.registerInfo = info
	n.initConfig()
	n.init()
	go n.heartbeat() // 心跳功能
}

// GetIp 获取单个ip
func (n *Nacos) GetIp() RegisterInfo {
	instance := n.SelOne(n.registerInfo.RemoteName)
	if instance == nil {
		return RegisterInfo{}
	}
	return RegisterInfo{Ip: instance.Ip, Port: uint16(instance.Port), ServiceName: instance.ServiceName}
}

// ListIp 获取ip
func (n *Nacos) ListIp() []RegisterInfo {
	instances := n.SelList(n.registerInfo.RemoteName)
	infos := make([]RegisterInfo, 0)
	if len(instances) == 0 {
		return infos
	}
	for _, item := range instances {
		infos = append(infos, RegisterInfo{Ip: item.Ip, Port: uint16(item.Port), ServiceName: item.ServiceName})
	}
	return infos
}

func NewNacos() *Nacos {
	return &Nacos{}
}

func (n *Nacos) initConfig() {
	// 创建serverConfig的另一种方式 -> 此处链接nacos的配置
	n.serverConfigs = []constant.ServerConfig{
		*constant.NewServerConfig(n.registerInfo.Ip, uint64(n.registerInfo.Port), constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}
	n.registerParam = vo.RegisterInstanceParam{
		Ip:          n.registerInfo.Ip,
		Port:        uint64(n.registerInfo.Port),
		ServiceName: n.registerInfo.ServiceName,
		GroupName:   "DEFAULT_GROUP",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	n.logoutParam = vo.DeregisterInstanceParam{
		Ip:          n.registerParam.Ip,
		Port:        n.registerParam.Port,
		ServiceName: n.registerParam.ServiceName,
	}
}

func (n *Nacos) init() {
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{ServerConfigs: n.serverConfigs})
	if err != nil {
		print(err)
	}
	n.namingClient = namingClient
	success, err := n.namingClient.RegisterInstance(n.registerParam)
	if err != nil || !success {
		print("注册失败:", err)
	}
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
//	log.Println("所有实例: ", instances)
//}

// SelList 只获取存活的实例
func (n *Nacos) SelList(serverName string) []model.Instance {
	// SelList 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
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

// SelOne 只获取存活的实例
func (n *Nacos) SelOne(serverName string) *model.Instance {
	// SelList 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		print(err)
	}
	return instances
}

// Heartbeat 心跳功能
func (n *Nacos) heartbeat() {
	for true {
		instance, err := n.namingClient.UpdateInstance(
			vo.UpdateInstanceParam{
				Ip:          n.registerParam.Ip,
				Port:        n.registerParam.Port,
				Weight:      n.registerParam.Weight,
				Enable:      n.registerParam.Enable,
				Healthy:     n.registerParam.Healthy,
				ServiceName: n.registerParam.ServiceName,
			})
		if err != nil || !instance {
			log.Panicln("更新实例失败,请检查Nacos!", err)
		}
		time.Sleep(1 * time.Second)
	}
}
